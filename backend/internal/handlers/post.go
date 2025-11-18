package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/middleware"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/server"
)

type PostHandler struct {
	postRepo         repository.PostRepository
	placeRepo        repository.PlaceRepository
	userRepo         repository.UserRepository
	commentRepo      repository.CommentRepository
	likeRepo         repository.LikeRepository
	notificationRepo repository.NotificationRepository
	validator        *validator.Validate
}

func NewPostHandler(
	postRepo repository.PostRepository,
	placeRepo repository.PlaceRepository,
	userRepo repository.UserRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	notificationRepo repository.NotificationRepository,
) *PostHandler {
	return &PostHandler{
		postRepo:         postRepo,
		placeRepo:        placeRepo,
		userRepo:         userRepo,
		commentRepo:      commentRepo,
		likeRepo:         likeRepo,
		notificationRepo: notificationRepo,
		validator:        validator.New(),
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	var req models.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err := h.placeRepo.GetByID(r.Context(), req.PlaceID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Place not found"})
		return
	}

	post := &models.Post{
		ID:          uuid.New(),
		UserID:      userID,
		PlaceID:     req.PlaceID,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.postRepo.Create(r.Context(), post); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
		return
	}

	// Create post images if provided
	if len(req.Images) > 0 {
		for i, imageURL := range req.Images {
			image := &models.PostImage{
				ID:        uuid.New(),
				PostID:    post.ID,
				ImageURL:  imageURL,
				SortOrder: i,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := h.postRepo.CreateImage(r.Context(), image); err != nil {
				// Log error but don't fail the post creation
				continue
			}
		}
	}

	// Get post with images for response
	postResponse := h.buildPostResponse(r.Context(), post)
	server.WriteJSON(w, http.StatusCreated, postResponse)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/posts/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	post, err := h.postRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	postResponse := h.buildPostResponse(r.Context(), post)
	server.WriteJSON(w, http.StatusOK, postResponse)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/posts/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	var req models.PostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Get existing post
	post, err := h.postRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	// TODO: Check if user owns the post

	// Update fields
	if req.Description != nil {
		post.Description = req.Description
	}
	post.UpdatedAt = time.Now()

	if err := h.postRepo.Update(r.Context(), post); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update post"})
		return
	}

	postResponse := h.buildPostResponse(r.Context(), post)
	server.WriteJSON(w, http.StatusOK, postResponse)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/posts/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	// TODO: Check if user owns the post

	if err := h.postRepo.Delete(r.Context(), id); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
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

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	posts, err := h.postRepo.ListFeed(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list posts"})
		return
	}

	responses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = h.buildPostResponse(r.Context(), post)
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"posts":  responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

func (h *PostHandler) ListPostsByUser(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/posts")]
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

	posts, err := h.postRepo.ListByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list user posts"})
		return
	}

	responses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = h.buildPostResponse(r.Context(), post)
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"posts":   responses,
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
		"count":   len(responses),
	})
}

func (h *PostHandler) ListPostsByPlace(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	placeIDStr := path[len("/api/places/") : len(path)-len("/posts")]
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
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

	posts, err := h.postRepo.ListByPlaceID(r.Context(), placeID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list place posts"})
		return
	}

	responses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = h.buildPostResponse(r.Context(), post)
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"posts":    responses,
		"place_id": placeID,
		"limit":    limit,
		"offset":   offset,
		"count":    len(responses),
	})
}

// buildPostResponse builds a complete post response with images, place, and user data
func (h *PostHandler) buildPostResponse(ctx context.Context, post *models.Post) models.PostResponse {
	response := post.ToResponse()

	// Get images
	images, err := h.postRepo.GetImagesByPostID(ctx, post.ID)
	if err == nil {
		response.Images = make([]models.PostImage, len(images))
		for i, img := range images {
			response.Images[i] = *img
		}
	}

	// Get place data
	place, err := h.placeRepo.GetByID(ctx, post.PlaceID)
	if err == nil {
		placeResponse := place.ToResponse()
		response.Place = &placeResponse
	}

	// Get user data
	user, err := h.userRepo.GetByID(ctx, post.UserID)
	if err == nil {
		userResponse := user.ToResponse()
		response.User = &userResponse
	}

	// Get engagement data
	if h.likeRepo != nil {
		if likes, err := h.likeRepo.CountPostLikes(ctx, post.ID); err == nil {
			response.LikesCount = likes
		}
		if userID, ok := middleware.GetUserIDFromContext(ctx); ok && userID != uuid.Nil {
			if liked, err := h.likeRepo.IsPostLikedByUser(ctx, post.ID, userID); err == nil {
				response.LikedByUser = liked
			}
		}
	}

	if h.commentRepo != nil {
		if comments, err := h.commentRepo.CountByPostID(ctx, post.ID); err == nil {
			response.CommentsCount = comments
		}
	}

	return response
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	postID, err := h.extractPostIDFromLikesPath(r.URL.Path)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	post, err := h.postRepo.GetByID(r.Context(), postID)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	if err := h.likeRepo.LikePost(r.Context(), postID, userID); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to like post"})
		return
	}

	if post.UserID != userID && h.notificationRepo != nil {
		data := map[string]string{
			"post_id": post.ID.String(),
		}
		notification := &models.Notification{
			UserID:  post.UserID,
			ActorID: userID,
			PostID:  &post.ID,
			Type:    models.NotificationTypeLikePost,
			Data:    data,
		}
		_ = h.createNotification(r.Context(), notification)
	}

	likes, _ := h.likeRepo.CountPostLikes(r.Context(), postID)
	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"post_id":     postID,
		"likes_count": likes,
		"liked":       true,
	})
}

func (h *PostHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	postID, err := h.extractPostIDFromLikesPath(r.URL.Path)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	post, err := h.postRepo.GetByID(r.Context(), postID)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	if err := h.likeRepo.UnlikePost(r.Context(), postID, userID); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to unlike post"})
		return
	}

	if h.notificationRepo != nil && post.UserID != userID {
		_ = h.notificationRepo.SoftDeleteByReference(
			r.Context(),
			post.UserID,
			userID,
			models.NotificationTypeLikePost,
			&post.ID,
			nil,
		)
	}

	likes, _ := h.likeRepo.CountPostLikes(r.Context(), postID)
	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"post_id":     postID,
		"likes_count": likes,
		"liked":       false,
	})
}

func (h *PostHandler) extractPostIDFromLikesPath(path string) (uuid.UUID, error) {
	const prefix = "/api/posts/"
	const suffix = "/likes"

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return uuid.Nil, fmt.Errorf("invalid likes path")
	}

	idStr := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	return uuid.Parse(idStr)
}

func (h *PostHandler) createNotification(ctx context.Context, notification *models.Notification) error {
	if h.notificationRepo == nil {
		return nil
	}

	notification.ID = uuid.New()
	if notification.Data == nil {
		notification.Data = map[string]string{}
	}

	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now

	return h.notificationRepo.Create(ctx, notification)
}
