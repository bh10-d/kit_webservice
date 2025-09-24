package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/db"
	"go-controller/dto"
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
	var req dto.CheckSiteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}
	
	payload := map[string]interface{}{
		"subDomain": req.SubDomain,
		"script":    "check_site.sh",
	}
	
	tag := req.Tag
	if tag == "" { 
		tag = "nginx" 
	}
	
	fmt.Println(payload)
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// CreateSite handles POST /create-site
func CreateSite(c *gin.Context) {
	var req dto.CreateSiteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}
	
	tag := req.Tag
	if tag == "" { 
		tag = "nginx" 
	}
	
	// First check if site exists
	checkPayload := map[string]interface{}{
		"subDomain": req.SubDomain,
		"script":    "check_site.sh",
	}
	checkResp, checkStatus := HandleFunc(tag, checkPayload)
	if checkStatus == 200 {
		// Site already exists, proceed with creation
		createPayload := map[string]interface{}{
			"subDomain": req.SubDomain,
			"script":    "create_site.sh",
		}
		result, status := HandleFunc(tag, createPayload)
		c.JSON(status, result)
	} else {
		c.JSON(checkStatus, checkResp)
	}
}

// UpdateSite handles PUT /update-site
func UpdateSite(c *gin.Context) {
	var req dto.UpdateSiteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}
	
	tag := req.Tag
	if tag == "" { 
		tag = "nginx" 
	}
	
	// Remove old site
	deletePayload := map[string]interface{}{
		"subDomain": req.OldSubDomain,
		"script":    "remove_site.sh",
	}
	deleteResp, deleteStatus := HandleFunc(tag, deletePayload)
	if deleteStatus != 200 {
		c.JSON(deleteStatus, deleteResp)
		return
	}
	
	// Check if new subdomain is available
	checkPayload := map[string]interface{}{
		"subDomain": req.NewSubDomain,
		"script":    "check_site.sh",
	}
	checkResp, checkStatus := HandleFunc(tag, checkPayload)
	if checkStatus != 200 {
		c.JSON(checkStatus, checkResp)
		return
	}
	
	// Create new site
	createPayload := map[string]interface{}{
		"subDomain": req.NewSubDomain,
		"script":    "create_site.sh",
	}
	result, status := HandleFunc(tag, createPayload)
	c.JSON(status, result)
}

// RemoveSite handles DELETE /remove-site
func RemoveSite(c *gin.Context) {
	var req dto.RemoveSiteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}
	
	tag := req.Tag
	if tag == "" { 
		tag = "nginx" 
	}
	
	payload := map[string]interface{}{
		"subDomain": req.SubDomain,
		"script":    "remove_site.sh",
	}
	
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// RegisterRunner handles POST /register
func RegisterRunner(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}
	
	runner := model.Runner{
		ID:       GenerateKey(),
		HostName: req.Name,
		IP:       req.IP,
		Tags:     req.Tags,
	}
	
	db.DB.Create(&runner)
	
	response := dto.RegisterResponse{
		Message: "Runner đăng ký thành công",
		Runner:  runner,
	}
	
	c.JSON(201, response)
}

