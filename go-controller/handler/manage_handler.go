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
func testSingleRunner(runner model.Runner) dto.HealthCheckResponse {
	start := time.Now()
	
	result := dto.HealthCheckResponse{
		RunnerID: runner.ID,
	}
	// Tạo test job payload
	testJobID := fmt.Sprintf("healthcheck_%s_%d", runner.ID, time.Now().Unix())
	testPayload := map[string]interface{}{
		"type":    "health_check",
		"command": "echo 'health_check_ok'",
	}

	

	// Gửi job đến runner
	msgID := rabbitmq.SendToQueueWithCustomID(runner.ID, testPayload, testJobID)
	if msgID == "" {
		result.Error = "Failed to send health check job to queue"
		result.ResponseTimeMs = int(time.Since(start).Milliseconds())
	}

	// Chờ response từ runner
	responseQueue := msgID + "_response"
	conn, err := rabbitmq.Dial()
	if err != nil {
		result.Error = fmt.Sprintf("RabbitMQ connection failed: %v", err)
		result.ResponseTimeMs = int(time.Since(start).Milliseconds())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		result.Error = fmt.Sprintf("RabbitMQ channel failed: %v", err)
		result.ResponseTimeMs = int(time.Since(start).Milliseconds())
	}
	defer ch.Close()

	// Declare response queue
	_, err = ch.QueueDeclare(responseQueue, true, false, false, false, nil)
	if err != nil {
		result.Error = fmt.Sprintf("Queue declare failed: %v", err)
		result.ResponseTimeMs = int(time.Since(start).Milliseconds())
	}

	timeout := 5
	waited := 0
	
	for waited < timeout {
		msg, ok, err := ch.Get(responseQueue, true)
		if err != nil {
			result.Error = fmt.Sprintf("Failed to get message: %v", err)
			break
		}
		
		if ok {
			var responseData map[string]interface{}
			if err := json.Unmarshal(msg.Body, &responseData); err == nil {
				if responseData["id"] == msgID {
					result.Alive = true
					result.ResponseTimeMs = int(time.Since(start).Milliseconds())
					result.ResponseID = responseData["id"].(string)
					result.Payload = testPayload
					result.Status = "done"
					return result
				}
			}
		}
		// time.Sleep(1 * time.Second)
		waited++
	}

	// Timeout - runner không response
	result.ResponseID = msgID
	result.Status = "error"
	result.Alive = false
	result.Error = "Health check timeout - runner did not respond"
	result.ResponseTimeMs = int(time.Since(start).Milliseconds())
	result.Payload = testPayload
	return result
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
	
	job := model.Job{
		RunnerID:        runnerID,
		MsgID:           fmt.Sprintf("%v", result.ResponseID), // an toàn, không panic
		Status:          fmt.Sprintf("%v", result.Status),
		RequestPayload:  toJSON(result.Payload),
		ResponsePayload: toJSON(result),
		Timeout:         false,
	}

	// Lưu job vào database
	if err := db.DB.Create(&job).Error; err != nil {
		fmt.Printf("Failed to create job record: %v\n", err)
	}

	var message	string

	// Cập nhật trạng thái runner
	if err := db.DB.Model(&runner).Where("id = ?", runnerID).Update("alive", result.Alive).Error; err != nil {
		fmt.Printf("Failed to update runner status: %v\n", err)
		message = fmt.Sprintf("Failed to update runner status: %v", err)
	}else{
		message = fmt.Sprintf("Runner %s is alive status updated", runnerID)
	}

	statusCode := 200
	if !result.Alive {
		// statusCode = 503
		statusCode = 200
	}


	c.JSON(statusCode, dto.ApiResponse{
		Status:  statusCode,
		Message: message,
		Data:    result,
	})
}