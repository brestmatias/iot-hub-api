package model

import "net/http"

type GenericJsonResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
	Status  int    `json:"status"`
}

func NewNotFoundResponse(message string) *GenericJsonResponse {
	return &GenericJsonResponse{
		Cause:   "not_found",
		Message: message,
		Status:  http.StatusNotFound,
	}
}
