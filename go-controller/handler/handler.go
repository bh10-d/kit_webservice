package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"sync"
	"github.com/gin-gonic/gin"
	"go-controller/db"
	"go-controller/model"
	"go-controller/rabbitmq"
	"go-controller/util"
)

func GenerateKey() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

var dbMu sync.Mutex

func FetchTags(tag string) []string {
	var runners []model.Runner
	dbMu.Lock()
	db.DB.Select("id", "tags").Where("tags LIKE ?", "%"+tag+"%").Find(&runners)
	dbMu.Unlock()
	var matchingIDs []string
	for _, r := range runners {
		tagList := util.SplitTags(r.Tags)
		for _, t := range tagList {
			if t == tag {
				matchingIDs = append(matchingIDs, r.ID)
			}
		}
	}
	return matchingIDs
}

func RunJob(runnerID string, payload map[string]interface{}, msgID string) (map[string]interface{}, int) {
	responseQueue := msgID + "_response"
	conn, err := rabbitmq.Dial()
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return map[string]interface{}{ "error": "RabbitMQ connection failed" }, 500
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return map[string]interface{}{ "error": "RabbitMQ channel failed" }, 500
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(responseQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare response queue: %v", err)
		return map[string]interface{}{ "error": "Queue declare failed" }, 500
	}
	waited := 0
	timeout := 10
	var response map[string]interface{}
	for waited < timeout {
		msg, ok, err := ch.Get(responseQueue, true)
		if err != nil {
			log.Printf("Failed to get message: %v", err)
			break
		}
		if ok {
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Body, &data); err == nil {
				if data["id"] == msgID {
					response = data
					break
				}
			}
		}
		time.Sleep(1 * time.Second)
		waited++
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if response != nil {
		job := model.Job{
			RunnerID: runnerID,
			MsgID: msgID,
			Status: fmt.Sprintf("%v", response["status"]),
			RequestPayload: toJSON(payload),
			ResponsePayload: toJSON(response),
			Timeout: false,
		}
		db.DB.Create(&job)
		if response["status"] == "error" {
			response["message"] = "Job failed on runner"
			return response, 422
		}
		response["message"] = "Done"
		return response, 200
	} else {
		job := model.Job{
			RunnerID: runnerID,
			MsgID: msgID,
			Status: "timeout",
			RequestPayload: toJSON(payload),
			Timeout: true,
		}
		db.DB.Create(&job)
		return map[string]interface{}{ "error": "⏳ Timeout chờ phản hồi từ runner" }, 504
	}
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func HandleFunc(tag string, payload map[string]interface{}) (gin.H, int) {
	runnerIDs := FetchTags(tag)
	if len(runnerIDs) == 0 {
		return gin.H{"error": fmt.Sprintf("No runner found with tag '%s'", tag)}, 404
	}
	results := []gin.H{}
	hasError := false
	for _, runnerID := range runnerIDs {
		msgID := rabbitmq.SendToQueue(runnerID, payload)
		result, statusCode := RunJob(runnerID, payload, msgID)
		if statusCode >= 400 {
			hasError = true
		}
		results = append(results, gin.H{
			"msg_id": msgID,
			"runner_id": runnerID,
			// "result": result,
			"log": result["log"],
			// "status_code": statusCode,
		})
	}
	response := gin.H{
		// "message": fmt.Sprintf("✅ Đã gửi đến %d runner", len(runnerIDs)),
		"message": fmt.Sprintf("Successful"),
		"data": results,
		"status": 200,
	}
	if hasError {
		// return response, 422
		response = gin.H{
			"message": fmt.Sprintf("Unsuccessful"),
			"data": []interface{}{},
			"status": 422,
		}
		return response, 422
	}
	return response, 207
}

// GetScripts handles GET /get-scripts
func GetScripts(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	c.JSON(200, gin.H{
		"scripts": util.ListScripts(),
	})
}

// CheckSite handles POST /check-site
func CheckSite(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	subDomain, ok := payload["subDomain"].(string)
	if !ok || subDomain == "" {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	payload["script"] = "check_site.sh"
	fmt.Println(payload)
	tag, _ := payload["tag"].(string)
	if tag == "" { tag = "nginx" }
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// CreateSite handles POST /create-site
func CreateSite(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	subDomain, ok := payload["subDomain"].(string)
	if !ok || subDomain == "" {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	tag, _ := payload["tag"].(string)
	if tag == "" { tag = "nginx" }
	payload["script"] = "check_site.sh"
	checkResp, checkStatus := HandleFunc(tag, payload)
	if checkStatus == 200 {
		payload["script"] = "create_site.sh"
		result, status := HandleFunc(tag, payload)
		c.JSON(status, result)
	} else {
		c.JSON(checkStatus, checkResp)
	}
}

// UpdateSite handles PUT /update-site
func UpdateSite(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	oldSubDomain, ok1 := payload["oldSubDomain"].(string)
	newSubDomain, ok2 := payload["newSubDomain"].(string)
	if !ok1 || !ok2 || oldSubDomain == "" || newSubDomain == "" {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	tag, _ := payload["tag"].(string)
	if tag == "" { tag = "nginx" }
	deletePayload := map[string]interface{}{
		"subDomain": oldSubDomain,
		"script": "remove_site.sh",
	}
	deleteResp, deleteStatus := HandleFunc(tag, deletePayload)
	if deleteStatus == 200 {
		c.JSON(deleteStatus, deleteResp)
		return
	}
	checkPayload := map[string]interface{}{
		"subDomain": newSubDomain,
	}
	checkResp, checkStatus := HandleFunc(tag, checkPayload)
	if checkStatus != 200 {
		c.JSON(checkStatus, checkResp)
		return
	}
	createPayload := map[string]interface{}{
		"subDomain": newSubDomain,
		"script": "create_site.sh",
	}
	result, status := HandleFunc(tag, createPayload)
	c.JSON(status, result)
}

// RemoveSite handles DELETE /remove-site
func RemoveSite(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	subDomain, ok := payload["subDomain"].(string)
	if !ok || subDomain == "" {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	tag, _ := payload["tag"].(string)
	if tag == "" { tag = "nginx" }
	payload["script"] = "remove_site.sh"
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// RegisterRunner handles POST /register
func RegisterRunner(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	name, _ := data["name"].(string)
	ip, _ := data["ip"].(string)
	tags, _ := data["tags"].(string)
	runner := model.Runner{
		ID: GenerateKey(),
		HostName: name,
		IP: ip,
		Tags: tags,
	}
	db.DB.Create(&runner)
	c.JSON(201, gin.H{
		"message": "Runner đăng ký thành công",
		"runner": runner,
	})
}

// GetJobs handles GET /get-jobs
func GetJobs(c *gin.Context) {
	var jobs []model.Job
	db.DB.Order("id desc").Limit(100).Find(&jobs)
	c.JSON(200, gin.H{
		"jobs": jobs,
	})
}
