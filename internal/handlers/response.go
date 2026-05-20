package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// sendErrorResponse sends a standardized error response
func sendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string, details ...string) {
	response := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}

	if len(details) > 0 {
		response.Details = details[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	render.JSON(w, r, response)
}

// sendSuccessResponse sends a standardized success response
func sendSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, data)
}
