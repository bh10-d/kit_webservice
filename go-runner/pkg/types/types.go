package types

import (
	"go-runner/internal/config"
)

type Runner struct {
	config *config.Config
}

type Job struct {
	ID        string `json:"id"`
	Script    string `json:"script"`
	// SubDomain string `json:"subdomain"`
	Parameters map[string]interface{} `json:"parameters"`
}

type Response struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Log     string `json:"log"`
}

type RegisterRequest struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Tags string `json:"tags"`
}

type RegisterResponse struct {
	Runner struct {
		ID string `json:"id"`
	} `json:"runner"`
}