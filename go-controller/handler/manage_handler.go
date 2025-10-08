package handler

import (
	// "fmt"
	"encoding/json"
	// "sync"
	"time"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/db"
	"go-controller/dto"
	"go-controller/model"
	"go-controller/service"
	"go-controller/rabbitmq"
	"fmt"
)

var manageScriptService = service.NewScriptService()

// GetJobs handles GET /get-jobs
func GetJobs(c *gin.Context) {
	var jobs []model.Job
	db.DB.Order("id desc").Limit(100).Find(&jobs)
	
	response := dto.ApiResponse{
		Status:  200,
		Message: "Jobs retrieved successfully",
		Data:    jobs,
	}
	c.JSON(200, response)
}

// func HealthCheck(c *gin.Context) {
// 	// Lấy danh sách tất cả runners
// 	var runners []model.Runner
// 	if err := db.DB.Find(&runners).Error; err != nil {
// 		c.JSON(500, dto.ApiResponse{
// 			Status:  500,
// 			Message: "Failed to fetch runners",
// 			Data: map[string]interface{}{
// 				"error": err.Error(),
// 			},
// 		})
// 		return
// 	}

// 	if len(runners) == 0 {
// 		c.JSON(200, dto.ApiResponse{
// 			Status:  200,
// 			Message: "No runners registered",
// 			Data: map[string]interface{}{
// 				"total_runners": 0,
// 				"alive_runners": 0,
// 				"dead_runners":  0,
// 				"runners":       []interface{}{},
// 			},
// 		})
// 		return
// 	}

// 	// Test tất cả runners song song
// 	results := testAllRunners(runners)
	
// 	aliveCount := 0
// 	deadCount := 0
	
// 	for _, result := range results {
// 		if result["status"] == "alive" {
// 			aliveCount++
// 		} else {
// 			deadCount++
// 		}
// 	}

// 	statusCode := 200
// 	message := "Health check completed"
	
// 	// Nếu quá nhiều runners chết thì báo warning
// 	if deadCount > 0 && deadCount >= len(runners)/2 {
// 		statusCode = 503
// 		message = "More than half of runners are down"
// 	}

// 	c.JSON(statusCode, dto.ApiResponse{
// 		Status:  statusCode,
// 		Message: message,
// 		Data: map[string]interface{}{
// 			"total_runners": len(runners),
// 			"alive_runners": aliveCount,
// 			"dead_runners":  deadCount,
// 			"runners":       results,
// 			"timestamp":     time.Now().Format(time.RFC3339),
// 		},
// 	})
// }

// // testAllRunners - Test tất cả runners song song
// func testAllRunners(runners []model.Runner) []map[string]interface{} {
// 	var wg sync.WaitGroup
// 	results := make([]map[string]interface{}, len(runners))
	
// 	for i, runner := range runners {
// 		wg.Add(1)
// 		go func(index int, r model.Runner) {
// 			defer wg.Done()
// 			results[index] = testSingleRunner(r)
// 		}(i, runner)
// 	}
	
// 	wg.Wait()
// 	return results
// }

// testSingleRunner - Test một runner cụ thể
func testSingleRunner(runner model.Runner) map[string]interface{} {
	start := time.Now()
	
	result := map[string]interface{}{
		"runner_id":  runner.ID,
		// "hostname":   runner.HostName,
		// "ip":         runner.IP,
		// "tags":       runner.Tags,
		"alive":     false,
		"response_time_ms": 0,
		// "error":      nil,
		// "last_seen":  runner.UpdatedAt.Format(time.RFC3339),
	}

	// Tạo test job payload
	testJobID := fmt.Sprintf("healthcheck_%s_%d", runner.ID, time.Now().Unix())
	testPayload := map[string]interface{}{
		"type":    "health_check",
		"command": "echo 'health_check_ok'",
		"timeout": 5, // 5 giây timeout
	}

	

	// Gửi job đến runner
	msgID := rabbitmq.SendToQueueWithCustomID(runner.ID, testPayload, testJobID)
	if msgID == "" {
		result["error"] = "Failed to send health check job to queue"
		result["response_time_ms"] = time.Since(start).Milliseconds()
		return result
	}

	// Chờ response từ runner
	responseQueue := msgID + "_response"
	conn, err := rabbitmq.Dial()
	if err != nil {
		result["error"] = fmt.Sprintf("RabbitMQ connection failed: %v", err)
		result["response_time_ms"] = time.Since(start).Milliseconds()
		return result
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		result["error"] = fmt.Sprintf("RabbitMQ channel failed: %v", err)
		result["response_time_ms"] = time.Since(start).Milliseconds()
		return result
	}
	defer ch.Close()

	// Declare response queue
	_, err = ch.QueueDeclare(responseQueue, true, false, false, false, nil)
	if err != nil {
		result["error"] = fmt.Sprintf("Queue declare failed: %v", err)
		result["response_time_ms"] = time.Since(start).Milliseconds()
		return result
	}

	// Chờ response trong 10 giây
	timeout := 5
	waited := 0
	
	for waited < timeout {
		msg, ok, err := ch.Get(responseQueue, true)
		if err != nil {
			result["error"] = fmt.Sprintf("Failed to get message: %v", err)
			break
		}
		
		if ok {
			var responseData map[string]interface{}
			if err := json.Unmarshal(msg.Body, &responseData); err == nil {
				if responseData["id"] == msgID {
					// Runner đã response!
					
					// c.JSON(200, response)
					result["alive"] = true
					result["response_time_ms"] = time.Since(start).Milliseconds()
					result["response_id"] = responseData["id"].(string)
					result["payload"] = testPayload
					result["status"] = "done"


					// response := dto.ApiResponse{
					// 	Status:  200,
					// 	Message: "Runners retrieved successfully",
					// 	Data:    result,
					// }
					return result
				}
			}
		}
		
		// time.Sleep(1 * time.Second)
		waited++
	}

	// Timeout - runner không response
	result["response_id"] = msgID
	result["status"] = "error"
	result["alive"] = false
	result["error"] = "Health check timeout - runner did not respond"
	result["response_time_ms"] = time.Since(start).Milliseconds()
	result["payload"] = testPayload
	return result
}

func GetRunners(c *gin.Context) {
	var runners []model.Runner
	db.DB.Limit(100).Find(&runners)
	
	response := dto.ApiResponse{
		Status:  200,
		Message: "Runners retrieved successfully",
		Data:    runners,
	}
	c.JSON(200, response)
}

func GetScripts(c *gin.Context) {
	scripts, err := manageScriptService.GetAllScripts()
	if err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "Failed to fetch scripts: " + err.Error()})
		return
	}
	
	response := dto.ScriptsResponse{
		Scripts: scripts,
	}
	c.JSON(200, response)
}

// GetScriptDetail handles GET /scripts/:id
func GetScriptDetail(c *gin.Context) {
	scriptID := c.Param("id")
	if scriptID == "" {
		c.JSON(400, dto.ErrorResponse{Error: "Script ID is required"})
		return
	}

	script, err := manageScriptService.GetScriptByID(scriptID)
	if err != nil {
		c.JSON(404, dto.ErrorResponse{Error: "Script not found: " + err.Error()})
		return
	}

	parameters, err := manageScriptService.GetScriptParameters(scriptID)
	if err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "Failed to get script parameters: " + err.Error()})
		return
	}

	response := dto.ScriptDetailResponse{
		Script:     *script,
		Parameters: parameters,
	}
	c.JSON(200, response)
}

func GetLogs(c *gin.Context) {
	var logs []model.Logs
	db.DB.Limit(100).Find(&logs)
	
	response := dto.ApiResponse{
		Status:  200,
		Message: "Logs retrieved successfully",
		Data:    logs,
	}
	c.JSON(200, response)
}

// GetJobGroup handles GET /jobs/:baseJobId - Get all jobs for a base job ID
func GetJobGroup(c *gin.Context) {
	baseJobID := c.Param("baseJobId")
	if baseJobID == "" {
		c.JSON(400, dto.ErrorResponse{Error: "Base Job ID is required"})
		return
	}

	summary, err := model.GetJobGroupSummary(db.DB, baseJobID)
	if err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "Failed to get job group: " + err.Error()})
		return
	}

	response := dto.ApiResponse{
		Status:  200,
		Message: "Job group retrieved successfully",
		Data:    summary,
	}
	c.JSON(200, response)
}




func CreateScript (c *gin.Context) {
	var script model.Scripts
	fmt.Printf("Received script creation request: %+v\n", script)
	if err := c.ShouldBindJSON(&script); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := manageScriptService.CreateScript(&script); err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "Failed to create script: " + err.Error()})
		return
	}

	// db.DB.Create(&script)

	response := dto.ApiResponse{
		Status:  201,
		Message: "Script created successfully",
		Data:    script,
	}
	c.JSON(201, response)
}
func UpdateScript (c *gin.Context) {
	var script model.Scripts
	id := c.Param("id")
	
	script.ScriptID = id
	// fmt.Println("Update ID from Param:", id)

	if err := c.ShouldBindJSON(&script); err != nil {
		fmt.Printf("Binding error: %v\n", err)
		c.JSON(400, dto.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// Debug: In ra dữ liệu nhận được
	fmt.Printf("Received script data: %+v\n", script)
	fmt.Printf("Param array: %v\n", script.Param)
	fmt.Printf("Tag array: %v\n", script.Tag)
	fmt.Printf("Runner array: %v\n", script.Runner)

	if err := manageScriptService.UpdateScript(&script); err != nil {
		fmt.Printf("Update error: %v\n", err)
		c.JSON(500, dto.ErrorResponse{Error: "Failed to update script: " + err.Error()})
		return
	}

	// Đọc lại từ database để kiểm tra
	updatedScript, err := manageScriptService.GetScriptByID(id)
	if err != nil {
		fmt.Printf("Error getting updated script: %v\n", err)
	} else {
		fmt.Printf("Updated script from DB: %+v\n", *updatedScript)
	}

	response := dto.ApiResponse{
		Status:  200,
		Message: "Script updated successfully",
		Data:    map[string]interface{}{
			"script":     updatedScript,
			"parameters": nil,
		},
	}
	c.JSON(200, response)
}


func UpdateScriptStatus (c *gin.Context) {
	// var script model.Scripts
	// if err := c.ShouldBindJSON(&script); err != nil {
	// 	c.JSON(400, dto.ErrorResponse{Error: "Invalid request payload"})
	// 	return
	// }

	// if err := manageScriptService.UpdateScriptStatus(&script); err != nil {
	// 	c.JSON(500, dto.ErrorResponse{Error: "Failed to update script status: " + err.Error()})
	// 	return
	// }

	id := c.Param("id")
	if id == "" {
		c.JSON(400, dto.ErrorResponse{Error: "Script ID is required"})
		return
	}

	if err := manageScriptService.UpdateScriptStatus(&model.Scripts{ScriptID: id, Status: true}); err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "Failed to update script status: " + err.Error()})
		return
	}

	response := dto.ApiResponse{
		Status:  200,
		Message: "Script status updated successfully",
	}
	c.JSON(200, response)
}

func DeleteScript (c *gin.Context) {
	// var script model.Scripts
	// if err := c.ShouldBindJSON(&script); err != nil {
	// 	c.JSON(400, dto.ErrorResponse{Error: "Invalid request payload"})
	// 	return
	// }

	// if err := manageScriptService.DeleteScript(script.ScriptID); err != nil {
	// 	c.JSON(500, dto.ErrorResponse{Error: "Failed to delete script: " + err.Error()})
	// 	return
	// }

	id := c.Param("id")
	if id == "" {
		c.JSON(400, dto.ErrorResponse{Error: "Script ID is required"})
		return
	}
	if err := manageScriptService.DeleteScript(id); err != nil {
		c.JSON(500, dto.ErrorResponse{Error: "Failed to delete script: " + err.Error()})
		return
	}
	
	response := dto.ApiResponse{
		Status:  200,
		Message: "Script status deleted successfully",
	}
	c.JSON(200, response)
}

// TestRunnerHealth - Test health của một runner cụ thể
func TestRunnerHealth(c *gin.Context) {
	runnerID := c.Param("id")
	if runnerID == "" {
		c.JSON(400, dto.ErrorResponse{Error: "Runner ID is required"})
		return
	}

	// Tìm runner trong database
	var runner model.Runner
	if err := db.DB.Where("id = ?", runnerID).First(&runner).Error; err != nil {
		c.JSON(404, dto.ErrorResponse{Error: "Runner not found"})
		return
	}

	// Test runner
	result := testSingleRunner(runner)
	
	statusCode := 200
	if !result["alive"].(bool) {
		statusCode = 503
	}

	// job := model.Job{
	// 	RunnerID:    runnerID,
	// 	MsgID:       result["id"].(string), // Lưu base job ID thay vì msgID có suffix
	// 	Status:      fmt.Sprintf("%v", result["status"]),
	// 	RequestPayload: toJSON(result["payload"]),
	// 	ResponsePayload: toJSON(result),
	// 	Timeout: false,
	// }

	job := model.Job{
		RunnerID:        runnerID,
		MsgID:           fmt.Sprintf("%v", result["response_id"]), // an toàn, không panic
		Status:          fmt.Sprintf("%v", result["status"]),
		RequestPayload:  toJSON(result["payload"]),
		ResponsePayload: toJSON(result),
		Timeout:         false,
	}


	db.DB.Create(&job)

	var message	string
	fmt.Println("Result alive:", result["alive"].(bool))
	if !result["alive"].(bool) {
		message = fmt.Sprintf("Runner %s is not alive", runnerID)
		db.DB.Model(&runner).
			Where("id = ?", runnerID).
			Update("alive", result["alive"].(bool))
		// c.JSON(400, dto.ErrorResponse{Error: message})

		c.JSON(200, dto.ApiResponse{
			Status:  statusCode,
			Message: message,
			Data:    result,
		})


		return
	} else {
		message = fmt.Sprintf("Runner %s is alive", runnerID)
		
		db.DB.Model(&runner).
			Where("id = ?", runnerID).
			Update("alive", result["alive"].(bool))
	}

	// fmt.Println(!result["alive"].(bool))

	// db.DB.Create(&job)

	// db.DB.Model(&runner).
	// 	Where("id = ?", runnerID).
	// 	Update("alive", result["alive"].(bool))

	

	// fmt.Println("RowsAffected:", tx.RowsAffected, "Error:", tx.Error)

	// db.DB.Model(&runner).Where("id = ?", runnerID).Update("alive", result["alive"].(bool))

	c.JSON(statusCode, dto.ApiResponse{
		Status:  statusCode,
		Message: message,
		Data:    result,
	})
}