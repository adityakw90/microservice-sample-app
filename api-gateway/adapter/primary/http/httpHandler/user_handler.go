package httpHandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	httpRequest "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/request"
	httpResponse "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/response"
	httpValidator "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/validator"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	userService service.UserService
	validator   *httpValidator.Validator
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   httpValidator.NewValidator(),
	}
}

// RegisterRoutes registers user routes on the given router.
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

// ListUsers handles GET /api/v1/users.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	param := params.ListUsersParam{
		ListParam: params.ListParam{
			PaginationParam: params.PaginationParam{
				Page:  parseIntQuery(r, "page", 1),
				Limit: parseIntQuery(r, "limit", 20),
			},
		},
	}

	if uids := r.URL.Query()["uid"]; len(uids) > 0 {
		param.UIDs = uids
	}
	if username := r.URL.Query().Get("username"); username != "" {
		param.Username = &username
	}
	if email := r.URL.Query().Get("email"); email != "" {
		param.Email = &email
	}
	if query := r.URL.Query().Get("query"); query != "" {
		param.Query = &query
	}
	if active := r.URL.Query().Get("active"); active != "" {
		activeBool := active == "true"
		param.Active = &activeBool
	}

	users, err := h.userService.ListUsers(ctx, &param)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UsersFromDomain(users))
}

// GetUser handles GET /api/v1/users/{uid}.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	user, err := h.userService.GetUser(ctx, uid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFromDomain(user))
}

// CreateUser handles POST /api/v1/users.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req httpRequest.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	createUserParam := params.CreateUserParam{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	uid, err := h.userService.CreateUser(ctx, &createUserParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"uid": uid})
}

// UpdateUser handles PUT/PATCH /api/v1/users/{uid}.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	var req httpRequest.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	updateUserParam := params.UpdateUserParam{
		UID:      uid,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Status:   req.Status,
	}

	if err := h.userService.UpdateUser(ctx, &updateUserParam); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.SuccessResponse{Success: true})
}

// DeleteUser handles DELETE /api/v1/users/{uid}.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	if err := h.userService.DeleteUser(ctx, uid); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.SuccessResponse{Success: true})
}

// GetUserProfile handles GET /api/v1/users/{uid}/profile.
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	profile, err := h.userService.GetUserProfile(ctx, uid)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.ProfileFromDomain(profile))
}

// UpdateUserProfile handles PUT/PATCH /api/v1/users/{uid}/profile.
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	var req httpRequest.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	updateProfileParam := params.UpdateUserProfileParam{
		UserUID:    uid,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Bio:        req.Bio,
		Attributes: req.Attributes,
	}

	if err := h.userService.UpdateUserProfile(ctx, &updateProfileParam); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.SuccessResponse{Success: true})
}

// ListUserDevices handles GET /api/v1/users/{uid}/devices.
func (h *UserHandler) ListUserDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]

	param := params.ListUserDevicesParam{
		ListParam: params.ListParam{
			PaginationParam: params.PaginationParam{
				Page:  parseIntQuery(r, "page", 1),
				Limit: parseIntQuery(r, "limit", 20),
			},
		},
		UserUID: uid,
	}

	devices, err := h.userService.ListUserDevices(ctx, &param)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.DevicesFromDomain(devices))
}

// RevokeUserDevice handles DELETE /api/v1/users/{uid}/devices/{deviceUid}.
func (h *UserHandler) RevokeUserDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	uid := vars["uid"]
	deviceUid := vars["deviceUid"]

	if err := h.userService.RevokeUserDevice(ctx, uid, deviceUid); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.SuccessResponse{Success: true})
}

// parseIntQuery parses an integer query parameter with a default value.
func parseIntQuery(r *http.Request, key string, defaultVal int) int {
	if val := r.URL.Query().Get(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
