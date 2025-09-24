package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/db"
	"go-controller/model"
)


// // GetScripts handles GET /get-scripts
// func GetScripts(c *gin.Context) {
// 	var payload map[string]interface{}
// 	if err := c.BindJSON(&payload); err != nil {
// 		c.JSON(400, gin.H{"error": "Missing JSON payload"})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"scripts": util.ListScripts(),
// 	})
// }

// CheckSite handles POST /check-site
func CheckSite(c *gin.Context) {
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
	fmt.Println(payload)
	tag, _ := payload["tag"].(string)
	if tag == "" { tag = "nginx" }
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// CreateSite handles POST /create-site
func CreateSite(c *gin.Context) {
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
	checkResp, checkStatus := HandleFunc(tag, payload)
	if checkStatus == 200 {
		payload["script"] = "create_site.sh"
		result, status := HandleFunc(tag, payload)
		c.JSON(status, result)
	} else {
		c.JSON(checkStatus, checkResp)
	}
}

// UpdateSite handles PUT /update-site
func UpdateSite(c *gin.Context) {
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
	deleteResp, deleteStatus := HandleFunc(tag, deletePayload)
	if deleteStatus == 200 {
		c.JSON(deleteStatus, deleteResp)
		return
	}
	checkPayload := map[string]interface{}{
		"subDomain": newSubDomain,
	}
	checkResp, checkStatus := HandleFunc(tag, checkPayload)
	if checkStatus != 200 {
		c.JSON(checkStatus, checkResp)
		return
	}
	createPayload := map[string]interface{}{
		"subDomain": newSubDomain,
		"script": "create_site.sh",
	}
	result, status := HandleFunc(tag, createPayload)
	c.JSON(status, result)
}

// RemoveSite handles DELETE /remove-site
func RemoveSite(c *gin.Context) {
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
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// RegisterRunner handles POST /register
func RegisterRunner(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Missing JSON payload"})
		return
	}
	name, _ := data["name"].(string)
	ip, _ := data["ip"].(string)
	tags, _ := data["tags"].(string)
	runner := model.Runner{
		ID: GenerateKey(),
		HostName: name,
		IP: ip,
		Tags: tags,
	}
	db.DB.Create(&runner)
	c.JSON(201, gin.H{
		"message": "Runner đăng ký thành công",
		"runner": runner,
	})
}

