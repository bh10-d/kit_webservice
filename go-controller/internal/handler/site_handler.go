package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	// "go-controller/util"
	"go-controller/internal/database"
	"go-controller/internal/dto"
	"go-controller/internal/model"
	"go-controller/internal/service"
	"go-controller/internal/util"
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

var scriptService = service.NewScriptService()

// ExecuteScript handles POST /execute-script - Generic script execution endpoint
func ExecuteScript(c *gin.Context) {
	var req dto.ScriptExecutionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}

	tag := req.Tag
	if tag == "" {
		tag = "nginx"
	}

	payload, err := scriptService.BuildScriptPayload(req.ScriptID, req.Parameters)
	if err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Script error: " + err.Error()})
		return
	}

	// fmt.Printf("Executing script %s with payload: %+v\n", req.ScriptID, payload)
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// CheckSite handles POST /check-site - Legacy endpoint
func CheckSite(c *gin.Context) {
	var req dto.CheckSiteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Invalid JSON payload or missing required fields"})
		return
	}

	// Convert to script execution request
	scriptReq, err := scriptService.ConvertSiteRequestToScriptRequest(req, "check")
	if err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Conversion error: " + err.Error()})
		return
	}


	// if err := util.CheckStatus(req.ScriptID); err != nil {
	// 	c.JSON(400, dto.ErrorResponse{Error: err.Error()})
	// 	return
	// }

	// ✅ Check status trong DB
	fmt.Printf("Checking status for script ID: %s\n", scriptReq.ScriptID)
	if err := util.CheckStatus(scriptReq.ScriptID); err != nil {
		// fmt.Println("Error checking script status:", scriptReq.ScriptID, err)
		c.JSON(400, dto.ErrorResponse{Error: scriptReq.ScriptID + ": " + err.Error()})
		return
	}

	tag := scriptReq.Tag
	if tag == "" {
		tag = "nginx"
	}

	payload, err := scriptService.BuildScriptPayload(scriptReq.ScriptID, scriptReq.Parameters)
	if err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Script error: " + err.Error()})
		return
	}

	fmt.Printf("Checking site with payload: %+v\n", payload)
	result, status := HandleFunc(tag, payload)
	c.JSON(status, result)
}

// CreateSite handles POST /create-site - Legacy endpoint
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

	// First check if site exists using script service
	checkReq, err := scriptService.ConvertSiteRequestToScriptRequest(req, "check")
	if err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Conversion error: " + err.Error()})
		return
	}

	checkPayload, err := scriptService.BuildScriptPayload(checkReq.ScriptID, checkReq.Parameters)
	if err != nil {
		c.JSON(400, dto.ErrorResponse{Error: "Script error: " + err.Error()})
		return
	}

	checkResp, checkStatus := HandleFunc(tag, checkPayload)
	if checkStatus == 200 {
		// Site check successful, proceed with creation
		createReq, err := scriptService.ConvertSiteRequestToScriptRequest(req, "create")
		if err != nil {
			c.JSON(400, dto.ErrorResponse{Error: "Conversion error: " + err.Error()})
			return
		}

		createPayload, err := scriptService.BuildScriptPayload(createReq.ScriptID, createReq.Parameters)
		if err != nil {
			c.JSON(400, dto.ErrorResponse{Error: "Script error: " + err.Error()})
			return
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

