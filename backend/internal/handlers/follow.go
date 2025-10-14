package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/middleware"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/server"
)

type FollowHandler struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
	validator  *validator.Validate
}

func NewFollowHandler(followRepo repository.FollowRepository, userRepo repository.UserRepository) *FollowHandler {
	return &FollowHandler{
		followRepo: followRepo,
		userRepo:   userRepo,
		validator:  validator.New(),
	}
}

func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	path := r.URL.Path
	// Extract user ID from /api/users/{id}/follow
	userIDStr := path[len("/api/users/") : len(path)-len("/follow")]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	if currentUserID == userID {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot follow yourself"})
		return
	}

	_, err = h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	exists, err := h.followRepo.IsFollowing(r.Context(), currentUserID, userID)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check follow status"})
		return
	}

	if exists {
		server.WriteJSON(w, http.StatusConflict, map[string]string{"error": "Already following this user"})
		return
	}

	follow := &models.Follow{
		ID:          uuid.New(),
		FollowerID:  currentUserID,
		FollowingID: userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.followRepo.CreateFollow(r.Context(), follow); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to follow user"})
		return
	}

	server.WriteJSON(w, http.StatusCreated, follow.ToResponse())
}

func (h *FollowHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/follow")]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	if err := h.followRepo.DeleteFollow(r.Context(), currentUserID, userID); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Follow relationship not found"})
		return
	}

	if err := h.followRepo.DeleteFollow(r.Context(), currentUserID, userID); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to unfollow user"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *FollowHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/following")]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
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

	users, err := h.followRepo.ListFollowing(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list following"})
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users":   responses,
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
		"count":   len(responses),
	})
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/followers")]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
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

	users, err := h.followRepo.ListFollowers(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list followers"})
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users":   responses,
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
		"count":   len(responses),
	})
}

func (h *FollowHandler) CheckFollowStatus(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/follow-status")]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Check if following
	isFollowing, err := h.followRepo.IsFollowing(r.Context(), currentUserID, userID)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check follow status"})
		return
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"is_following": isFollowing,
		"user_id":      userID,
		"follower_id":  currentUserID,
	})
}

func (h *FollowHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/stats")]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	_, err = h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Get user stats
	stats, err := h.followRepo.GetUserStats(r.Context(), userID)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get user stats"})
		return
	}

	server.WriteJSON(w, http.StatusOK, stats)
}
