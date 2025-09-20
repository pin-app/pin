package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/server"
)

type UserHandler struct {
	userRepo  repository.UserRepository
	validator *validator.Validate
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user := &models.User{
		ID:          uuid.New(),
		Email:       req.Email,
		Username:    req.Username,
		Bio:         req.Bio,
		Location:    req.Location,
		DisplayName: req.DisplayName,
		PfpURL:      req.PfpURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		return
	}

	server.WriteJSON(w, http.StatusCreated, user.ToResponse())
}

// GetUser handles GET /api/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	server.WriteJSON(w, http.StatusOK, user.ToResponse())
}

// UpdateUser handles PUT /api/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Get existing user
	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Update fields
	if req.Username != nil {
		user.Username = req.Username
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Location != nil {
		user.Location = req.Location
	}
	if req.DisplayName != nil {
		user.DisplayName = req.DisplayName
	}
	if req.PfpURL != nil {
		user.PfpURL = req.PfpURL
	}
	user.UpdatedAt = time.Now()

	if err := h.userRepo.Update(r.Context(), user); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		return
	}

	server.WriteJSON(w, http.StatusOK, user.ToResponse())
}

// DeleteUser handles DELETE /api/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	if err := h.userRepo.Delete(r.Context(), id); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

// ListUsers handles GET /api/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userRepo.List(r.Context(), limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list users"})
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users":  responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

// SearchUsers handles GET /api/users/search
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Query parameter 'q' is required"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userRepo.Search(r.Context(), query, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to search users"})
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users":  responses,
		"query":  query,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}
