package dto

// import "go-controller/model"

type SiteOperationRequest struct {
	SubDomain string `json:"subDomain" binding:"required"`
	ScriptID  string `json:"scriptId" binding:"required"`
	Tag       string `json:"tag,omitempty"`
}

type UpdateSiteRequest struct {
	OldSubDomain string `json:"oldSubDomain" binding:"required"`
	NewSubDomain string `json:"newSubDomain" binding:"required"`
	Tag          string `json:"tag,omitempty"`
}

type SiteOperationResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Output  string `json:"output,omitempty"`
}

