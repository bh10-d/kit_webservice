package dto

// import "go-controller/model"

// Legacy DTOs for backward compatibility
type SiteRequest struct {
	SubDomain string `json:"subDomain" binding:"required"`
	Tag       string `json:"tag,omitempty"`
	ScriptID  string    `json:"script_id" binding:"required"`
} 

type CheckSiteRequest = SiteRequest
type CreateSiteRequest = SiteRequest
type RemoveSiteRequest = SiteRequest

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Sort     string `form:"sort" json:"sort"`
	Order    string `form:"order" json:"order"`
	Search   string `form:"search" json:"search"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Status     int            `json:"status"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}