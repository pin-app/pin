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

type PostHandler struct {
	postRepo  repository.PostRepository
	placeRepo repository.PlaceRepository
	userRepo  repository.UserRepository
	validator *validator.Validate
}

func NewPostHandler(postRepo repository.PostRepository, placeRepo repository.PlaceRepository, userRepo repository.UserRepository) *PostHandler {
	return &PostHandler{
		postRepo:  postRepo,
		placeRepo: placeRepo,
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// CreatePost handles POST /api/posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req models.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Verify place exists
	_, err := h.placeRepo.GetByID(r.Context(), req.PlaceID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Place not found"})
		return
	}

	post := &models.Post{
		ID:          uuid.New(),
		UserID:      uuid.New(), // TODO: Get from session/auth
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

// GetPost handles GET /api/posts/{id}
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

// UpdatePost handles PUT /api/posts/{id}
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

// DeletePost handles DELETE /api/posts/{id}
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

// ListPosts handles GET /api/posts
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.postRepo.ListFeed(r.Context(), uuid.New(), limit, offset) // TODO: Get user ID from session
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

// ListPostsByUser handles GET /api/users/{id}/posts
func (h *PostHandler) ListPostsByUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path like "/api/users/{id}/posts"
	path := r.URL.Path
	userIDStr := path[len("/api/users/") : len(path)-len("/posts")]
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

// ListPostsByPlace handles GET /api/places/{id}/posts
func (h *PostHandler) ListPostsByPlace(w http.ResponseWriter, r *http.Request) {
	// Extract place ID from path like "/api/places/{id}/posts"
	path := r.URL.Path
	placeIDStr := path[len("/api/places/") : len(path)-len("/posts")]
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
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

	return response
}
