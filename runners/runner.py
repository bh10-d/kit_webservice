import pika
import json
import time
import subprocess
import os
import socket
import json
import requests
from dotenv import load_dotenv, set_key

# Tự động load từ file .env
ENV_PATH = ".env"
load_dotenv(dotenv_path=ENV_PATH)
# load_dotenv()

KEY_ID = os.getenv("KEY_ID")
SERVER_URL = os.getenv("SERVER_URL", "http://localhost:5000")
RUNNER_TAGS = os.getenv("RUNNER_TAGS")
MESSAGE_QUEUE = os.getenv("MESSAGE_QUEUE")
SCRIPTS_PATH = os.getenv("SCRIPTS_PATH")
hostname = socket.gethostname()
IPAddr = socket.gethostbyname(hostname)

def register_runner():
    global KEY_ID

    if KEY_ID:
        print(f"✅ Runner đã đăng ký. Bỏ qua đăng ký lại.", flush=True)
        return

    print("🛰 Bắt đầu đăng ký runner...", flush=True)

    payload = {
        "name": hostname,
        "ip": IPAddr,  # Bạn có thể tự lấy IP nội bộ nếu cần
        "tags": RUNNER_TAGS,
    }

    for i in range(10):
        try:
            res = requests.post(f"{SERVER_URL}/register", json=payload, timeout=3)
            if res.status_code == 201:
                data = res.json()
                runner_token = data['runner']['id']
                print("✅ Đăng ký thành công!", flush=True)
                print("Runner Token:", runner_token, flush=True)

                # Ghi vào .env bằng set_key
                set_key(ENV_PATH, "KEY_ID", runner_token)

                # Ghi lại vào biến toàn cục
                KEY_ID = runner_token

                print(f"✅ KEY_ID đã được ghi vào {ENV_PATH}", flush=True)
                return
            else:
                print(f"⚠️ Server trả về lỗi ({res.status_code}): {res.text}", flush=True)
        except requests.exceptions.RequestException as e:
            print(f"⚠️ Kết nối thất bại (lần {i+1}/10): {e}", flush=True)
        time.sleep(2)

    print("❌ Không thể đăng ký runner sau nhiều lần thử", flush=True)
    exit(1)

# ========= Main Login ===========

def log(log):
    print(log, flush=True)


def wait_for_rabbitmq(max_retries=10, delay=3):
    for i in range(max_retries):
        try:
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(host=MESSAGE_QUEUE)
            )
            log("✅ Kết nối thành công đến RabbitMQ!")
            log(f"keyid: {KEY_ID}")
            return connection
        except pika.exceptions.AMQPConnectionError as e:
            log(f"⚠️ Kết nối thất bại (lần {i+1}/{max_retries}): {e}")
            time.sleep(delay)
    print("❌ Không thể kết nối đến RabbitMQ sau nhiều lần thử.", flush=True)
    exit(1)


def callback(ch, method, properties, body):

    log("📥 Đã gọi callback.")
    log(f"📦 Raw body: {body}")
    try:
        data = json.loads(body)
        log(f"📨 Parsed message: {data}")
        log(f"📨 Nhận được message: {data}")
        # script = f"./check.sh"
        script = data.get("script")
        subDomain = data.get("subDomain")
        log(f"📁 Script được lấy ra: {script}")
        # script_result = execute_script(script)
        success, output = execute_script(script, subDomain)

        response = {
            "id": data.get("id"),
            "status": "done" if success else "error",
            "message": f"Script '{script}' executed",
            "log": output
        }

        send_response(KEY_ID, response)

    except Exception as e:
        log(f"❌ Lỗi khi xử lý message: {e}")
        log(f"Nội dung thô: {body}")
        response = {
            "id": data.get("id"),
            "status": "error",
            "message": str(e) or "Unknown error"
        }
        send_response(KEY_ID, response)

def send_response(queue, data):
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host=MESSAGE_QUEUE)
    )
    channel = connection.channel()

    # Tạo tên queue phản hồi mới
    response_queue = f"{queue}_response"
    channel.queue_declare(queue=response_queue, durable=True)

    channel.basic_publish(
        exchange='',
        routing_key=response_queue,
        body=json.dumps(data)
    )

    log(f"📨 Đã gửi phản hồi vào queue: {response_queue}")
    connection.close()

def execute_script(script, subDomain):
    log(f"Thực thi Script: {script}")
    try:
        process = subprocess.run([f"{SCRIPTS_PATH}/{script}",f"{subDomain}"], capture_output=True, text=True, check=True)
        log(f"✅ Output:\n{process.stdout}")
        return True, process.stdout.strip()
    except subprocess.CalledProcessError as e:
        log(f"❌ Script lỗi:\n{e.stdout}\n{e.stderr}")
        return False, (e.stdout + e.stderr).strip()


if __name__ == '__main__':
    register_runner()

    # 🔄 Reload .env sau khi ghi KEY_ID mới
    load_dotenv(override=True)
    KEY_ID = os.getenv("KEY_ID")

    connection = wait_for_rabbitmq()
    channel = connection.channel()

    channel.queue_declare(queue=KEY_ID, durable=True)
    channel.basic_consume(queue=KEY_ID, on_message_callback=callback, auto_ack=True)

    log("👂 Đang chờ message...")
    channel.start_consuming()