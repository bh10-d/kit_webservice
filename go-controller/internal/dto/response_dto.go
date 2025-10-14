package dto

// import "go-controller/model"

type ApiResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

