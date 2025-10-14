package dto

import "go-controller/internal/model"

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
	IP   string `json:"ip" binding:"required"`
	Tags string `json:"tags"`
}

type RegisterResponse struct {
	Message string       `json:"message"`
	Runner  model.Runner `json:"runner"`
}

type HealthCheckResponse struct {
	RunnerID     string      `json:"runner_id"`
	ResponseID   string      `json:"response_id"`
	Alive        bool        `json:"alive"`
	Error        string      `json:"error"`
	ResponseTimeMs int       `json:"response_time_ms"`
	Payload      interface{} `json:"payload"`
	Status       string      `json:"status"`
}

