package httpHandler

import (
	"encoding/json"
	"net/http"

	httpResponse "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/response"
)

// respondJSON sends a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, httpResponse.ErrorResponse{Error: message})
}

// respondValidationError sends a validation error response.
func respondValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]map[string]string{"errors": errors})
}
