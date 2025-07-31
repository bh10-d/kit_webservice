
package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"go-controller/db"
	"go-controller/handler"
	"go-controller/model"
)

func main() {
	db.InitDB()
	r := gin.Default()

	r.POST("/check-site", func(c *gin.Context) {
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
		tag, _ := payload["tag"].(string)
		if tag == "" { tag = "nginx" }
		result, status := handler.HandleFunc(tag, payload)
		c.JSON(status, result)
	})

	r.POST("/create-site", func(c *gin.Context) {
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
		checkResp, checkStatus := handler.HandleFunc(tag, payload)
		if checkStatus == 200 {
			payload["script"] = "create_site.sh"
			result, status := handler.HandleFunc(tag, payload)
			c.JSON(status, result)
		} else {
			c.JSON(checkStatus, checkResp)
		}
	})

	r.PUT("/update-site", func(c *gin.Context) {
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
		deleteResp, deleteStatus := handler.HandleFunc(tag, deletePayload)
		if deleteStatus == 200 {
			c.JSON(deleteStatus, deleteResp)
			return
		}
		checkPayload := map[string]interface{}{
			"subDomain": newSubDomain,
		}
		checkResp, checkStatus := handler.HandleFunc(tag, checkPayload)
		if checkStatus != 200 {
			c.JSON(checkStatus, checkResp)
			return
		}
		createPayload := map[string]interface{}{
			"subDomain": newSubDomain,
			"script": "create_site.sh",
		}
		result, status := handler.HandleFunc(tag, createPayload)
		c.JSON(status, result)
	})

	r.DELETE("/remove-site", func(c *gin.Context) {
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
		result, status := handler.HandleFunc(tag, payload)
		c.JSON(status, result)
	})

	r.POST("/register", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": "Missing JSON payload"})
			return
		}
		name, _ := data["name"].(string)
		ip, _ := data["ip"].(string)
		tags, _ := data["tags"].(string)
		runner := model.Runner{
			ID: handler.GenerateKey(),
			Name: name,
			IP: ip,
			Tags: tags,
		}
		db.DB.Create(&runner)
		c.JSON(201, gin.H{
			"message": "Runner đăng ký thành công",
			"runner": runner,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}
