package dto

import "go-controller/internal/model"

type ScriptExecutionRequest struct {
	ScriptID   string                 `json:"scriptId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Parameters map[string]interface{} `json:"parameters" binding:"required" example:"{\"subDomain\":\"example\"}"`
	Tag        string                 `json:"tag,omitempty" example:"nginx"`
}

type ScriptDetailResponse struct {
	Script      model.Scripts            `json:"script"`
	Parameters  []ScriptParameterInfo    `json:"parameters"`
}

type ScriptParameterInfo struct {
	Name        string `json:"name" example:"subDomain"`
	Type        string `json:"type" example:"string"`
	Required    bool   `json:"required" example:"true"`
	Description string `json:"description" example:"The subdomain name"`
}

type ScriptsResponse struct {
	Scripts []model.Scripts `json:"scripts"`
}

