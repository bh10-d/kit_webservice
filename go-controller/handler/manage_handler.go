package handler

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/db"
	"go-controller/dto"
	"go-controller/model"
	"go-controller/service"
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
// func UpdateScript (c *gin.Context) {}
// func DeleteScript (c *gin.Context) {
// 	var script model.Scripts
// 	if err := c.ShouldBindJSON(&script); err != nil {
// 		c.JSON(400, dto.ErrorResponse{Error: "Invalid request payload"})
// 		return
// 	}

// 	db.DB.Delete(&script)

// 	response := dto.ApiResponse{
// 		Status:  200,
// 		Message: "Script deleted successfully",
// 		Data:    script,
// 	}
// 	c.JSON(200, response)
// }