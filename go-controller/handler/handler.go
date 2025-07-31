
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
			"runner_id": runnerID,
			"msg_id": msgID,
			"result": result,
			"status_code": statusCode,
		})
	}
	response := gin.H{
		"message": fmt.Sprintf("✅ Đã gửi đến %d runner", len(runnerIDs)),
		"results": results,
	}
	if hasError {
		return response, 422
	}
	return response, 207
}
