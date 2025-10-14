
package main

import (
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"go-controller/internal/database"
	"go-controller/internal/routes"
	// "go-controller/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Không tìm thấy file .env, dùng env của hệ thống")
	}
	
	// Initialize database
	db.InitDB()
	
	// // Initialize default scripts
	// scriptService := service.NewScriptService()
	// if err := scriptService.InitializeDefaultScripts(); err != nil {
	// 	log.Printf("⚠️ Failed to initialize default scripts: %v", err)
	// } else {
	// 	log.Println("✅ Default scripts initialized successfully")
	// }
	
	// Initialize Gin router
	r := gin.Default()
	r.Use(cors.Default())
	
	// Setup all routes
	routes.SetupRoutes(r)
	
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}
