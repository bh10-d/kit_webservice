package routes

import (
	"github.com/gin-gonic/gin"
	"go-controller/handler"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
	// Generic script execution endpoint
	r.POST("/execute-script", handler.ExecuteScript)
	
	// Scripts management routes
	r.GET("/scripts", handler.GetScripts)
	r.GET("/scripts/:id", handler.GetScriptDetail)
	r.POST("/scripts", handler.CreateScript)
	// r.PUT("/scripts/:id", handler.UpdateScript)
	// r.DELETE("/scripts/:id", handler.DeleteScript)
	
	// Legacy site management routes (for backward compatibility)
	r.POST("/check-site", handler.CheckSite)
	r.POST("/create-site", handler.CreateSite)
	r.PUT("/update-site", handler.UpdateSite)
	r.DELETE("/remove-site", handler.RemoveSite)
	
	// Runner management routes
	r.POST("/register", handler.RegisterRunner)
	
	// Management routes
	r.GET("/get-jobs", handler.GetJobs)
	r.GET("/jobs/:baseJobId", handler.GetJobGroup)  // New: Get jobs by base job ID
	r.GET("/get-runners", handler.GetRunners)
	r.GET("/get-scripts", handler.GetScripts) // Legacy endpoint
	r.GET("/get-logs", handler.GetLogs)
}
