package handler

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/db"
	"go-controller/model"
)

// GetJobs handles GET /get-jobs
func GetJobs(c *gin.Context) {
	var jobs []model.Job
	db.DB.Order("id desc").Limit(100).Find(&jobs)
	c.JSON(200, gin.H{
		"jobs": jobs,
	})
}

func GetRunners(c *gin.Context) {
	var runners []model.Runner
	db.DB.Limit(100).Find(&runners)
	c.JSON(200, gin.H{
		"runners": runners,
	})
}

func GetScripts(c *gin.Context) {
	var scripts []model.Scripts
	db.DB.Limit(100).Find(&scripts)
	c.JSON(200, gin.H{
		"scripts": scripts,
	})
}


func GetLogs(c *gin.Context) {
	var logs []model.Logs
	db.DB.Limit(100).Find(&logs)
	c.JSON(200, gin.H{
		"logs": logs,
	})
}