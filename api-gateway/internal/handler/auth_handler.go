package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/service"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRoutes registers auth routes on the given router
func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	// Auth routes
	r.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", h.RefreshToken).Methods(http.MethodPost)
	r.HandleFunc("/auth/validate", h.ValidateToken).Methods(http.MethodPost)
	r.HandleFunc("/auth/verify-pin", h.VerifyPin).Methods(http.MethodPost)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.authService.Login(ctx, req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ValidateToken handles POST /api/v1/auth/validate
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.authService.ValidateToken(ctx, req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// VerifyPin handles POST /api/v1/auth/verify-pin
func (h *AuthHandler) VerifyPin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.VerifyPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.authService.VerifyPin(ctx, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
