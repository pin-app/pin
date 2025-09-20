package handlers

import (
	"context"
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

type CommentHandler struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	validator   *validator.Validate
}

func NewCommentHandler(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository) *CommentHandler {
	return &CommentHandler{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		validator:   validator.New(),
	}
}

// CreateComment handles POST /api/comments
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req models.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Verify post exists
	_, err := h.postRepo.GetByID(r.Context(), req.PostID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Post not found"})
		return
	}

	// If parent_id is provided, verify parent comment exists
	if req.ParentID != nil {
		_, err := h.commentRepo.GetByID(r.Context(), *req.ParentID)
		if err != nil {
			server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Parent comment not found"})
			return
		}
	}

	comment := &models.Comment{
		ID:        uuid.New(),
		PostID:    req.PostID,
		UserID:    uuid.New(), // TODO: Get from session/auth
		ParentID:  req.ParentID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.commentRepo.Create(r.Context(), comment); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create comment"})
		return
	}

	commentResponse := h.buildCommentResponse(r.Context(), comment)
	server.WriteJSON(w, http.StatusCreated, commentResponse)
}

// GetComment handles GET /api/comments/{id}
func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/comments/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.commentRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		return
	}

	commentResponse := h.buildCommentResponse(r.Context(), comment)
	server.WriteJSON(w, http.StatusOK, commentResponse)
}

// UpdateComment handles PUT /api/comments/{id}
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/comments/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
		return
	}

	var req models.CommentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Get existing comment
	comment, err := h.commentRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		return
	}

	// TODO: Check if user owns the comment

	// Update fields
	comment.Content = req.Content
	comment.UpdatedAt = time.Now()

	if err := h.commentRepo.Update(r.Context(), comment); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update comment"})
		return
	}

	commentResponse := h.buildCommentResponse(r.Context(), comment)
	server.WriteJSON(w, http.StatusOK, commentResponse)
}

// DeleteComment handles DELETE /api/comments/{id}
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/comments/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
		return
	}

	// TODO: Check if user owns the comment

	if err := h.commentRepo.Delete(r.Context(), id); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

// ListCommentsByPost handles GET /api/posts/{id}/comments
func (h *CommentHandler) ListCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Path[len("/api/posts/") : len("/api/posts/")+36] // UUID length
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
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

	comments, err := h.commentRepo.ListByPostID(r.Context(), postID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list comments"})
		return
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = h.buildCommentResponse(r.Context(), comment)
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"comments": responses,
		"post_id":  postID,
		"limit":    limit,
		"offset":   offset,
		"count":    len(responses),
	})
}

// ListCommentsByUser handles GET /api/users/{id}/comments
func (h *CommentHandler) ListCommentsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Path[len("/api/users/") : len("/api/users/")+36] // UUID length
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

	comments, err := h.commentRepo.ListByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list user comments"})
		return
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = h.buildCommentResponse(r.Context(), comment)
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"comments": responses,
		"user_id":  userID,
		"limit":    limit,
		"offset":   offset,
		"count":    len(responses),
	})
}

// GetCommentReplies handles GET /api/comments/{id}/replies
func (h *CommentHandler) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/comments/") : len("/api/comments/")+36] // UUID length
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
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

	comments, err := h.commentRepo.GetReplies(r.Context(), id, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get comment replies"})
		return
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = h.buildCommentResponse(r.Context(), comment)
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"comments":  responses,
		"parent_id": id,
		"limit":     limit,
		"offset":    offset,
		"count":     len(responses),
	})
}

// builds a complete comment response with user data
func (h *CommentHandler) buildCommentResponse(ctx context.Context, comment *models.Comment) models.CommentResponse {
	response := comment.ToResponse()

	// Get user data
	user, err := h.userRepo.GetByID(ctx, comment.UserID)
	if err == nil {
		userResponse := user.ToResponse()
		response.User = &userResponse
	}

	return response
}
