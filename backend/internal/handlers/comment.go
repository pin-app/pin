package handlers

import (
	"context"
	"encoding/json"
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

type CommentHandler struct {
	commentRepo      repository.CommentRepository
	postRepo         repository.PostRepository
	userRepo         repository.UserRepository
	notificationRepo repository.NotificationRepository
	validator        *validator.Validate
}

func NewCommentHandler(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository, notificationRepo repository.NotificationRepository) *CommentHandler {
	return &CommentHandler{
		commentRepo:      commentRepo,
		postRepo:         postRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		validator:        validator.New(),
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	var req models.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	post, err := h.postRepo.GetByID(r.Context(), req.PostID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Post not found"})
		return
	}

	var parentComment *models.Comment
	if req.ParentID != nil {
		parentComment, err = h.commentRepo.GetByID(r.Context(), *req.ParentID)
		if err != nil {
			server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Parent comment not found"})
			return
		}
	}

	comment := &models.Comment{
		ID:        uuid.New(),
		PostID:    req.PostID,
		UserID:    userID,
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

	h.dispatchCommentNotifications(r.Context(), comment, post, parentComment)
	server.WriteJSON(w, http.StatusCreated, commentResponse)
}

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

	comment, err := h.commentRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		return
	}

	comment.Content = req.Content
	comment.UpdatedAt = time.Now()

	if err := h.commentRepo.Update(r.Context(), comment); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update comment"})
		return
	}

	commentResponse := h.buildCommentResponse(r.Context(), comment)
	server.WriteJSON(w, http.StatusOK, commentResponse)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/comments/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
		return
	}

	if err := h.commentRepo.Delete(r.Context(), id); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Comment not found"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *CommentHandler) ListCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Path[len("/api/posts/") : len("/api/posts/")+36] // UUID length
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

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

func (h *CommentHandler) ListCommentsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Path[len("/api/users/") : len("/api/users/")+36] // UUID length
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

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

func (h *CommentHandler) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/comments/") : len("/api/comments/")+36] // UUID length
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

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

func (h *CommentHandler) dispatchCommentNotifications(ctx context.Context, comment *models.Comment, post *models.Post, parentComment *models.Comment) {
	if h.notificationRepo == nil {
		return
	}

	if post.UserID != comment.UserID {
		data := map[string]string{
			"post_id":    post.ID.String(),
			"comment_id": comment.ID.String(),
		}
		h.createNotification(ctx, &models.Notification{
			UserID:    post.UserID,
			ActorID:   comment.UserID,
			PostID:    &post.ID,
			CommentID: &comment.ID,
			Type:      models.NotificationTypeCommentPost,
			Data:      data,
		})
	}

	if parentComment != nil && parentComment.UserID != comment.UserID && parentComment.UserID != post.UserID {
		data := map[string]string{
			"post_id":           post.ID.String(),
			"comment_id":        comment.ID.String(),
			"parent_comment_id": parentComment.ID.String(),
		}
		h.createNotification(ctx, &models.Notification{
			UserID:    parentComment.UserID,
			ActorID:   comment.UserID,
			PostID:    &post.ID,
			CommentID: &comment.ID,
			Type:      models.NotificationTypeCommentReply,
			Data:      data,
		})
	}
}

func (h *CommentHandler) createNotification(ctx context.Context, notification *models.Notification) {
	notification.ID = uuid.New()
	if notification.Data == nil {
		notification.Data = map[string]string{}
	}
	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now
	_ = h.notificationRepo.Create(ctx, notification)
}
