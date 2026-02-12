package httpHandler

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/params"
	httpRequest "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/request"
	httpResponse "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/response"
	httpValidator "github.com/adityakw90/microservice-sample-app/api-gateway/adapter/primary/http/validator"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/port/service"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
)

// UserFileHandler handles HTTP requests for file operations.
type UserFileHandler struct {
	fileService service.UserFileService
	validator   *httpValidator.Validator
}

// NewUserFileHandler creates a new file handler.
func NewUserFileHandler(fileService service.UserFileService) *UserFileHandler {
	return &UserFileHandler{
		fileService: fileService,
		validator:   httpValidator.NewValidator(),
	}
}

// RegisterRoutes registers file routes on the given router.
func (h *UserFileHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/files", h.ListFiles).Methods(http.MethodGet)
	r.HandleFunc("/files", h.UploadFile).Methods(http.MethodPost)
	r.HandleFunc("/files/{uid}", h.GetFile).Methods(http.MethodGet)
	r.HandleFunc("/files/{uid}", h.UpdateFile).Methods(http.MethodPut, http.MethodPatch)
	r.HandleFunc("/files/{uid}", h.DeleteFile).Methods(http.MethodDelete)
}

// ListFiles handles GET /api/v1/files.
func (h *UserFileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 20
	}

	listParam := params.ListUserFilesParam{
		Page:  page,
		Limit: limit,
	}

	if userUID := r.URL.Query().Get("user_uid"); userUID != "" {
		listParam.UserUIDs = []string{userUID}
	}
	if fileType := r.URL.Query().Get("file_type"); fileType != "" {
		listParam.FileType = &fileType
	}

	files, err := h.fileService.ListFiles(ctx, &listParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFilesFromDomain(files))
}

// UploadFile handles POST /api/v1/files.
func (h *UserFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		respondError(w, http.StatusRequestEntityTooLarge, "File too large (max 10MB)")
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Get user_uid from form or JWT
	userUID := r.FormValue("user_uid")
	if userUID == "" {
		// TODO: Get from JWT context
		respondError(w, http.StatusBadRequest, "user_uid required")
		return
	}

	// Get visibility (default private)
	visibility := int32(0) // private
	if visStr := r.FormValue("visibility"); visStr == "1" || visStr == "true" {
		visibility = 1 // public
	}

	uploadParam := params.UploadUserFileParam{
		UserUID:    userUID,
		FileName:   handler.Filename,
		FileData:   fileData,
		MimeType:   handler.Header.Get("Content-Type"),
		Visibility: &visibility,
	}

	newFile, err := h.fileService.UploadFile(ctx, &uploadParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, httpResponse.UserFileFromDomain(newFile))
}

// GetFile handles GET /api/v1/files/{uid}.
func (h *UserFileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		respondError(w, http.StatusBadRequest, "File UID required")
		return
	}

	file, err := h.fileService.GetFile(ctx, uid)
	if err != nil {
		respondError(w, http.StatusNotFound, "File not found")
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFileFromDomain(file))
}

// UpdateFile handles PUT/PATCH /api/v1/files/{uid}.
func (h *UserFileHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		respondError(w, http.StatusBadRequest, "File UID required")
		return
	}

	var req httpRequest.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := h.validator.ValidateStruct(req); len(errors) > 0 {
		respondValidationError(w, errors)
		return
	}

	updateParam := params.UpdateUserFileParam{
		UID:        uid,
		FileName:   req.FileName,
		Visibility: req.Visibility,
	}

	updatedFile, err := h.fileService.UpdateFile(ctx, &updateParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.UserFileFromDomain(updatedFile))
}

// DeleteFile handles DELETE /api/v1/files/{uid}.
func (h *UserFileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	uid := vars["uid"]

	if uid == "" {
		respondError(w, http.StatusBadRequest, "File UID required")
		return
	}

	if err := h.fileService.DeleteFile(ctx, uid); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, httpResponse.SuccessResponse{Success: true})
}
