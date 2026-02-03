package handler

// Common HTTP response types

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ValidationResponse represents a validation error response
type ValidationResponse struct {
	Error  string              `json:"error"`
	Fields map[string]string   `json:"fields,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{Error: message}
}

// NewErrorResponseWithCode creates a new error response with a code
func NewErrorResponseWithCode(message, code string) *ErrorResponse {
	return &ErrorResponse{Error: message, Code: code}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{Success: true}
}

// NewSuccessResponseWithMessage creates a new success response with a message
func NewSuccessResponseWithMessage(message string) *SuccessResponse {
	return &SuccessResponse{Success: true, Message: message}
}
