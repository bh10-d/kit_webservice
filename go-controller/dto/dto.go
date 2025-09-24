// package dto

// import "go-controller/model"

package dto

import "go-controller/model"

// SiteRequest represents the common request payload for site operations
type SiteRequest struct {
	SubDomain string `json:"subDomain" binding:"required" example:"example"`
	Tag       string `json:"tag,omitempty" example:"nginx"`
} // @name SiteRequest

// CheckSiteRequest represents the request payload for checking a site
type CheckSiteRequest = SiteRequest // @name CheckSiteRequest

// CreateSiteRequest represents the request payload for creating a site
type CreateSiteRequest = SiteRequest // @name CreateSiteRequest

// RemoveSiteRequest represents the request payload for removing a site
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
	Scripts []string `json:"scripts" example:"check_site.sh,create_site.sh,remove_site.sh"`
} // @name ScriptsResponse

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
