package handler

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/db"
	"go-controller/dto"
	"go-controller/model"
)

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
	var scripts []model.Scripts
	db.DB.Limit(100).Find(&scripts)
	
	response := dto.ApiResponse{
		Status:  200,
		Message: "Scripts retrieved successfully",
		Data:    scripts,
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