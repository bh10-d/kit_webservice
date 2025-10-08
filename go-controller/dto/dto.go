package dto

import "go-controller/model"

// ScriptExecutionRequest represents a generic request for executing scripts
type ScriptExecutionRequest struct {
	ScriptID   string                 `json:"scriptId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Parameters map[string]interface{} `json:"parameters" binding:"required" example:"{\"subDomain\":\"example\"}"`
	Tag        string                 `json:"tag,omitempty" example:"nginx"`
} // @name ScriptExecutionRequest

// SiteOperationRequest represents request for site-specific operations
type SiteOperationRequest struct {
	SubDomain string `json:"subDomain" binding:"required" example:"example"`
	ScriptID  string `json:"scriptId" binding:"required" example:"check-site-script"`
	Tag       string `json:"tag,omitempty" example:"nginx"`
} // @name SiteOperationRequest

// Legacy DTOs for backward compatibility
type SiteRequest struct {
	SubDomain string `json:"subDomain" binding:"required" example:"example"`
	Tag       string `json:"tag,omitempty" example:"nginx"`
	ScriptID  string    `json:"script_id" binding:"required" example:"1"`
} // @name SiteRequest

type CheckSiteRequest = SiteRequest // @name CheckSiteRequest
type CreateSiteRequest = SiteRequest // @name CreateSiteRequest  
type RemoveSiteRequest = SiteRequest // @name RemoveSiteRequest

// UpdateSiteRequest represents the request payload for updating a site
type UpdateSiteRequest struct {
	OldSubDomain string `json:"oldSubDomain" binding:"required" example:"old-example"`
	NewSubDomain string `json:"newSubDomain" binding:"required" example:"new-example"`
	Tag          string `json:"tag,omitempty" example:"nginx"`
} // @name UpdateSiteRequest

// RegisterRequest represents the request payload for registering a runner
type RegisterRequest struct {
	Name string `json:"name" binding:"required" example:"web-server-01"`
	IP   string `json:"ip" binding:"required" example:"192.168.1.100"`
	Tags string `json:"tags" example:"nginx,web"`
} // @name RegisterRequest

// ScriptsResponse represents the response for getting scripts
type ScriptsResponse struct {
	Scripts []model.Scripts `json:"scripts"`
} // @name ScriptsResponse

// ScriptDetailResponse represents detailed script information
type ScriptDetailResponse struct {
	Script      model.Scripts            `json:"script"`
	Parameters  []ScriptParameterInfo    `json:"parameters"`
} // @name ScriptDetailResponse

// ScriptParameterInfo represents information about script parameters
type ScriptParameterInfo struct {
	Name        string `json:"name" example:"subDomain"`
	Type        string `json:"type" example:"string"`
	Required    bool   `json:"required" example:"true"`
	Description string `json:"description" example:"The subdomain name"`
} // @name ScriptParameterInfo

// ApiResponse represents a generic API response
type ApiResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error message"`
} // @name ApiResponse

// SiteOperationResponse represents response for site operations
type SiteOperationResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Operation completed successfully"`
	Output  string `json:"output,omitempty" example:"Site operation output"`
} // @name SiteOperationResponse

// RegisterResponse represents response for runner registration
type RegisterResponse struct {
	Message string       `json:"message" example:"Runner đăng ký thành công"`
	Runner  model.Runner `json:"runner"`
} // @name RegisterResponse

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Missing JSON payload"`
} // @name ErrorResponse


// type HealthCheckResponse struct {
// 	runner_id   string
// 	response_id string
// 	alive       bool
// 	error       bool
// 	response_time int
// 	payload     interface{}
// 	status      string
// }


type HealthCheckResponse struct {
	RunnerID     string      `json:"runner_id"`
	ResponseID   string      `json:"response_id"`
	Alive        bool        `json:"alive"`
	Error        string      `json:"error"`
	ResponseTimeMs int      `json:"response_time_ms"`
	Payload      interface{} `json:"payload"`
	Status       string      `json:"status"`
}