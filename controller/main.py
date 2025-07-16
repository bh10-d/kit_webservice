from flask import Flask, request, jsonify
import pika
import json
import uuid
import sys
import time
import threading

from models import db, Runner, Job#, Technique #TaskResult, PendingKey
import secrets


app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///runners.db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

db.init_app(app)

with app.app_context():
    db.create_all()

RABBITMQ_HOST = '192.168.248.1'

registered_runners = {}

# Bộ nhớ tạm để lưu phản hồi từ runner
response_lock = threading.Lock()

def send_to_queue(runner_id, data):
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host=RABBITMQ_HOST)
    )
    channel = connection.channel()

    channel.queue_declare(queue=runner_id, durable=True)

    if 'id' not in data:
        data['id'] = str(uuid.uuid4())

    data['reply_to'] = runner_id

    channel.basic_publish(
        exchange='',
        routing_key=runner_id,
        body=json.dumps(data)
    )

    connection.close()
    return data['id']

def fetchTags(tag="nginx"):
    search_tag = f"%{tag}%"
    runners = db.session.query(Runner.id, Runner.tags).filter(Runner.tags.like(search_tag)).all()

    matching_ids = []
    for runner_id, tags in runners:
        tag_list = [t.strip() for t in tags.split(",")]
        if tag in tag_list:
            matching_ids.append(runner_id)

    return matching_ids

def generate_key():
    key_id = secrets.token_hex(16)
    return key_id

def runJob(runner_id, payload, msg_id):
    # response_queue = f"{runner_id}_response"
    response_queue = f"{msg_id}_response"

    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host=RABBITMQ_HOST)
    )
    channel = connection.channel()
    channel.queue_declare(queue=response_queue, durable=True)

    waited = 0
    timeout = 10  
    response = None

    print(f"🧪 Đang lắng nghe queue: {response_queue}, chờ msg_id: {msg_id}", flush=True)
    
    while waited < timeout:
        method_frame, header_frame, body = channel.basic_get(queue=response_queue, auto_ack=True)
        if method_frame:
            try:
                data = json.loads(body)
                if data.get("id") == msg_id:
                    response = data
                    break
            except Exception as e:
                print("❌ Lỗi parse response:", e)
        time.sleep(1)
        waited += 1

    connection.close()

    if response:
        try:
            with response_lock:
                job = Job(
                    # id=msg_id,
                    runner_id=runner_id,
                    msg_id=msg_id,
                    status=response.get("status"),
                    request_payload=json.dumps(payload),
                    response_payload=json.dumps(response),
                    timeout=False,
                )
                db.session.add(job)
                db.session.commit()
        except Exception as e:
            print(f"❌ Lỗi khi lưu kết quả vào DB: {e}", flush=True)

        # ✅ Kiểm tra lỗi logic từ runner
        if response.get("status") == "error":
            return {
                # 'result': response
                **response,  # unpack từng key ra ngoài
                'message': 'Job failed on runner',
            }, 422 

        return {
            # 'result': response
            **response,  # unpack từng key ra ngoài
            'message': 'Done',
        }, 200

    else:
        # Lưu Job bị timeout
        try:
            with response_lock:
                job = Job(
                    # id=msg_id,
                    runner_id=runner_id,
                    msg_id=msg_id,
                    status='timeout',
                    request_payload=json.dumps(payload),
                    timeout=True
                )
                db.session.add(job)
                db.session.commit()
        except Exception as e:
            print(f"❌ Lỗi khi lưu kết quả timeout vào DB: {e}", flush=True)

        return {
            'error': '⏳ Timeout chờ phản hồi từ runner'
        }, 504


def handleFunc(tag, payload):
    runner_ids = fetchTags(tag)
    if not runner_ids:
        return jsonify({'error': f"No runner found with tag '{tag}'"}), 404

    results = []
    has_error = False

    for runner_id in runner_ids:
        msg_id = send_to_queue(runner_id, payload)
        result, status_code = runJob(runner_id, payload, msg_id)

        # Ghi nhận nếu có bất kỳ lỗi nào
        if status_code >= 400:
            has_error = True

        results.append({
            "runner_id": runner_id,
            "msg_id": msg_id,
            "result": result,
            "status_code": status_code
        })

    response = {
        "message": f"✅ Đã gửi đến {len(runner_ids)} runner",
        "results": results
    }

    # Nếu có lỗi ở bất kỳ runner nào → trả về 422
    return jsonify(response), 422 if has_error else 207

def check_site_logic(payload, subDomain):
    print(f"🔍 Check subdomain: {subDomain}", flush=True)

    payload["script"] = "check_site.sh"
    tag = payload.get("tag", "nginx")
    runner_ids = fetchTags(tag)

    if not runner_ids:
        return jsonify({'error': f"No runner found with tag '{tag}'"}), 404

    runner_id = runner_ids[0]  # Chọn 1 runner đầu tiên để check
    msg_id = send_to_queue(runner_id, payload)
    result, status_code = runJob(runner_id, payload, msg_id)

    return result, status_code

@app.route('/check-site', methods=['POST'])
def checkSite():
    payload = request.get_json()
    subDomain = payload.get("subDomain")

    if not payload:
        return jsonify({'error': 'Missing JSON payload'}), 400
    
    if not subDomain:
        return jsonify({'error': 'Missing JSON payload'}), 400
    
    return check_site_logic(payload, subDomain)

@app.route('/create-site', methods=['POST'])
def createSite():
    payload = request.json
    subDomain = payload.get("subDomain")
    if not payload or not subDomain:
        return jsonify({'error': 'Missing JSON payload'}), 400

    tag = payload.get("tag", "nginx")

    # Gọi check logic trước
    check_resp, check_status = check_site_logic(payload, subDomain)

    if check_status == 200:
        payload["script"] = "create_site.sh"
        return handleFunc(tag, payload)
    else:
        return check_resp, check_status


@app.route('/update-site', methods=['PUT'])
def updateSite():
    payload = request.get_json()
    oldSubDomain = payload.get("oldSubDomain")
    newSubDomain = payload.get("newSubDomain")

    if not payload or not oldSubDomain or not newSubDomain:
        return jsonify({'error': 'Missing JSON payload'}), 400

    tag = payload.get("tag", "nginx")

    ## Bước 1: Xóa site cũ
    delete_payload = {
        "subDomain": oldSubDomain,
        "script": "remove_site.sh",
    }
    delete_resp, delete_status = handleFunc(tag, delete_payload)
    if delete_status == 200:
        return delete_resp, delete_status

    ## Bước 2: Check site mới
    check_resp, check_status = check_site_logic({
        "subDomain": newSubDomain,
    }, newSubDomain)

    if check_status != 200:
        return check_resp, check_status

    ## Bước 3: Tạo site mới
    create_payload = {
        "subDomain": newSubDomain,
        "script": "create_site.sh",
    }
    return handleFunc(tag, create_payload)

@app.route('/remove-site', methods=['DELETE'])
def removeSite():
    payload = request.get_json()
    subDomain = payload.get("subDomain")
    if not payload:
        return jsonify({'error': 'Missing JSON payload'}), 400
    
    if not subDomain:
        return jsonify({'error': 'Missing JSON payload'}), 400
    
    tag = payload.get("tag", "nginx")

    # Gọi check logic trước
    check_resp, check_status = check_site_logic(payload, subDomain)

    if check_status >= 400:
        payload["script"] = "remove_site.sh"
        return handleFunc(tag, payload)
    else:
        return check_resp, check_status

@app.route('/retry', methods=['POST'])
def retry_failed_job():
    data = request.json
    msg_id = data.get("msg_id")  # ⚠️ giờ lấy msg_id từ request

    if not msg_id:
        return jsonify({"error": "Missing msg_id"}), 400

    job = Job.query.filter_by(msg_id=msg_id).first()

    if not job:
        return jsonify({"error": "Job not found"}), 404

    print(job, flush=True)
    print("id:", job.runner_id, flush=True)

    runner_id = job.runner_id
    statusJob = job.status

    if statusJob == 'error':
        # Gửi lại job
        new_msg_id = send_to_queue(runner_id, json.loads(job.request_payload))
        result, status_code = runJob(runner_id, json.loads(job.request_payload), new_msg_id)

        return jsonify({
            "retry_status": "sent",
            "msg_id": new_msg_id,
            "result": result,
            "code": status_code
        }), status_code
    else:
        return jsonify({
            "msg_id": msg_id,
            "result": "Retry fail. Job was successful"
        }), 422


@app.route('/register', methods=['POST'])
# @app.route('/register', methods=['GET'])
def register_runner():
    try:
        data = request.get_json()
        name = data.get("name")
        ip = data.get("ip")
        tags = data.get("tags")

        # technique_name = Technique.query.filter_by(technique_name="nginx").first() #dang lam
        # technique_path = Technique.query.filter_by(technique_path="nginx").first()
        # script_path =  Technique.query.filter_by(technique_path="nginx").first()

        runner = Runner(
            id=generate_key(),
            name=name,
            ip=ip,
            tags=tags,
        )
        db.session.add(runner)
        db.session.commit()

        return jsonify({
            "message": "Runner đăng ký thành công",
            "runner": runner.to_dict()
        }), 201

    except Exception as e:
        print("Lỗi:", e, flush=True)
        return jsonify({"error": "Lỗi server"}), 500


if __name__ == '__main__':
    # fetchTags()
    print("📋 Registered routes:", flush=True)
    for rule in app.url_map.iter_rules():
        print(rule, flush=True)

    app.run(host='0.0.0.0', port=5000)
