package httpHandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	httpRequest "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/request"
	httpResponse "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/response"
	httpValidator "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/validator"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

// AuthHandler handles HTTP requests for authentication operations.
type AuthHandler struct {
	authService service.AuthService
	validator    *httpValidator.Validator
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:    httpValidator.NewValidator(),
	}
}

// RegisterRoutes registers auth routes on the given router.
func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/refresh", h.RefreshToken).Methods(http.MethodPost)
	r.HandleFunc("/auth/validate", h.ValidateToken).Methods(http.MethodPost)
	r.HandleFunc("/auth/verify-pin", h.VerifyPin).Methods(http.MethodPost)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req httpRequest.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	loginParam := params.LoginParam{
		Identifier:        req.Identifier,
		IdentifierType:    req.IdentifierType,
		Password:          req.Password,
		DeviceFingerprint: req.DeviceFingerprint,
		DeviceName:        req.DeviceName,
	}

	tokens, err := h.authService.Login(ctx, &loginParam)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.TokensFromDomain(tokens))
}

// RefreshToken handles POST /api/v1/auth/refresh.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req httpRequest.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	refreshParam := params.RefreshTokenParam{
		RefreshToken: req.RefreshToken,
	}

	tokens, err := h.authService.RefreshToken(ctx, &refreshParam)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.TokensFromDomain(tokens))
}

// ValidateToken handles POST /api/v1/auth/validate.
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req httpRequest.ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	validateParam := params.ValidateTokenParam{
		AccessToken: req.AccessToken,
	}

	claims, err := h.authService.ValidateToken(ctx, &validateParam)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.TokenClaimsFromDomain(claims))
}

// VerifyPin handles POST /api/v1/auth/verify-pin.
func (h *AuthHandler) VerifyPin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req httpRequest.VerifyPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	verifyPinParam := params.VerifyPinParam{
		UID:  req.UID,
		Code: req.Code,
	}

	valid, err := h.authService.VerifyPin(ctx, &verifyPinParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"valid": valid})
}
