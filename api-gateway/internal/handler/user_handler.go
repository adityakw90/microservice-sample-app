package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/service"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers user routes on the given router
func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	// User routes
	r.HandleFunc("/users", h.ListUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{uid}", h.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{uid}", h.UpdateUser).Methods(http.MethodPut, http.MethodPatch)
	r.HandleFunc("/users/{uid}", h.DeleteUser).Methods(http.MethodDelete)

	// User profile routes
	r.HandleFunc("/users/{uid}/profile", h.GetUserProfile).Methods(http.MethodGet)
	r.HandleFunc("/users/{uid}/profile", h.UpdateUserProfile).Methods(http.MethodPut, http.MethodPatch)

	// User device routes
	r.HandleFunc("/users/{uid}/devices", h.ListUserDevices).Methods(http.MethodGet)
	r.HandleFunc("/users/{uid}/devices/{deviceUid}", h.RevokeUserDevice).Methods(http.MethodDelete)
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	req := service.ListUsersRequest{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
	}

	// Optional filters
	if uids := r.URL.Query()["uid"]; len(uids) > 0 {
		req.UIDs = uids
	}
	if username := r.URL.Query().Get("username"); username != "" {
		req.Username = &username
	}
	if email := r.URL.Query().Get("email"); email != "" {
		req.Email = &email
	}
	if query := r.URL.Query().Get("query"); query != "" {
		req.Query = &query
	}
	if active := r.URL.Query().Get("active"); active != "" {
		activeBool := active == "true"
		req.Active = &activeBool
	}

	resp, err := h.userService.ListUsers(ctx, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// GetUser handles GET /api/v1/users/{uid}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	user, err := h.userService.GetUser(ctx, uid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	uid, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"uid": uid})
}

// UpdateUser handles PUT/PATCH /api/v1/users/{uid}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	var req service.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.UID = uid

	if err := h.userService.UpdateUser(ctx, req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// DeleteUser handles DELETE /api/v1/users/{uid}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	if err := h.userService.DeleteUser(ctx, uid); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// GetUserProfile handles GET /api/v1/users/{uid}/profile
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	profile, err := h.userService.GetUserProfile(ctx, uid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, profile)
}

// UpdateUserProfile handles PUT/PATCH /api/v1/users/{uid}/profile
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	var req service.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.UserUID = uid

	if err := h.userService.UpdateUserProfile(ctx, req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// ListUserDevices handles GET /api/v1/users/{uid}/devices
func (h *UserHandler) ListUserDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	req := service.ListUserDevicesRequest{
		UserUID: uid,
		Page:    parseIntQuery(r, "page", 1),
		Limit:   parseIntQuery(r, "limit", 20),
	}

	resp, err := h.userService.ListUserDevices(ctx, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// RevokeUserDevice handles DELETE /api/v1/users/{uid}/devices/{deviceUid}
func (h *UserHandler) RevokeUserDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]
	deviceUid := vars["deviceUid"]

	if err := h.userService.RevokeUserDevice(ctx, uid, deviceUid); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Helper functions

func parseIntQuery(r *http.Request, key string, defaultVal int) int {
	if val := r.URL.Query().Get(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
