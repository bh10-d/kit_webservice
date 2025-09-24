package routes

import (
	"github.com/gin-gonic/gin"
	"go-controller/handler"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
	// Scripts routes
	// r.POST("/get-scripts", handler.GetScripts)
	
	// Site management routes
	r.POST("/check-site", handler.CheckSite)
	r.POST("/create-site", handler.CreateSite)
	r.PUT("/update-site", handler.UpdateSite)
	r.DELETE("/remove-site", handler.RemoveSite)
	
	// Runner management routes
	r.POST("/register", handler.RegisterRunner)
	
	// Management routes
	r.GET("/get-jobs", handler.GetJobs)
	r.GET("/get-runners", handler.GetRunners)
	r.GET("/get-scripts", handler.GetScripts)
	r.GET("/get-logs", handler.GetLogs)
}
