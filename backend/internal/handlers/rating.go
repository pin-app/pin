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

type RatingHandler struct {
	ratingRepo repository.RatingRepository
	placeRepo  repository.PlaceRepository
	userRepo   repository.UserRepository
	validator  *validator.Validate
}

func NewRatingHandler(ratingRepo repository.RatingRepository, placeRepo repository.PlaceRepository, userRepo repository.UserRepository) *RatingHandler {
	return &RatingHandler{
		ratingRepo: ratingRepo,
		placeRepo:  placeRepo,
		userRepo:   userRepo,
		validator:  validator.New(),
	}
}

func (h *RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	placeIDStr := r.URL.Path[len("/api/places/") : len("/api/places/")+36] // UUID length
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	var req models.PlaceRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err = h.placeRepo.GetByID(r.Context(), placeID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Place not found"})
		return
	}

	userID := uuid.New() // TODO: Get from session/auth

	rating := &models.PlaceRating{
		ID:        uuid.New(),
		UserID:    userID,
		PlaceID:   placeID,
		Rating:    req.Rating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.ratingRepo.CreateRating(r.Context(), rating); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create rating"})
		return
	}

	server.WriteJSON(w, http.StatusCreated, rating)
}

func (h *RatingHandler) GetRating(w http.ResponseWriter, r *http.Request) {
	placeIDStr := r.URL.Path[len("/api/places/") : len("/api/places/")+36] // UUID length
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	userID := uuid.New() // TODO: Get from session/auth

	rating, err := h.ratingRepo.GetRating(r.Context(), userID, placeID)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Rating not found"})
		return
	}

	server.WriteJSON(w, http.StatusOK, rating)
}

func (h *RatingHandler) UpdateRating(w http.ResponseWriter, r *http.Request) {
	placeIDStr := r.URL.Path[len("/api/places/") : len("/api/places/")+36] // UUID length
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	var req models.PlaceRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID := uuid.New() // TODO: Get from session/auth

	rating := &models.PlaceRating{
		UserID:    userID,
		PlaceID:   placeID,
		Rating:    req.Rating,
		UpdatedAt: time.Now(),
	}

	if err := h.ratingRepo.UpdateRating(r.Context(), rating); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Rating not found"})
		return
	}

	server.WriteJSON(w, http.StatusOK, rating)
}

func (h *RatingHandler) DeleteRating(w http.ResponseWriter, r *http.Request) {
	placeIDStr := r.URL.Path[len("/api/places/") : len("/api/places/")+36] // UUID length
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	userID := uuid.New() // TODO: Get from session/auth

	if err := h.ratingRepo.DeleteRating(r.Context(), userID, placeID); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Rating not found"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *RatingHandler) ListRatingsByPlace(w http.ResponseWriter, r *http.Request) {
	placeIDStr := r.URL.Path[len("/api/places/") : len("/api/places/")+36] // UUID length
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

	ratings, err := h.ratingRepo.GetRatingsByPlaceID(r.Context(), placeID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list ratings"})
		return
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"ratings":  ratings,
		"place_id": placeID,
		"limit":    limit,
		"offset":   offset,
		"count":    len(ratings),
	})
}

func (h *RatingHandler) GetAverageRating(w http.ResponseWriter, r *http.Request) {
	placeIDStr := r.URL.Path[len("/api/places/") : len("/api/places/")+36] // UUID length
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	avgRating, count, err := h.ratingRepo.GetAverageRating(r.Context(), placeID)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get average rating"})
		return
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"average_rating": avgRating,
		"count":          count,
		"place_id":       placeID,
	})
}

func (h *RatingHandler) CreateComparison(w http.ResponseWriter, r *http.Request) {
	var req models.PlaceComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err := h.placeRepo.GetByID(r.Context(), req.BetterPlaceID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Better place not found"})
		return
	}

	_, err = h.placeRepo.GetByID(r.Context(), req.WorsePlaceID)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Worse place not found"})
		return
	}

	userID := uuid.New() // TODO: Get from session/auth

	comparison := &models.PlaceComparison{
		ID:            uuid.New(),
		UserID:        userID,
		BetterPlaceID: req.BetterPlaceID,
		WorsePlaceID:  req.WorsePlaceID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.ratingRepo.CreateComparison(r.Context(), comparison); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create comparison"})
		return
	}

	server.WriteJSON(w, http.StatusCreated, comparison)
}

func (h *RatingHandler) ListComparisonsByUser(w http.ResponseWriter, r *http.Request) {
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

	comparisons, err := h.ratingRepo.GetComparisonsByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list comparisons"})
		return
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"comparisons": comparisons,
		"user_id":     userID,
		"limit":       limit,
		"offset":      offset,
		"count":       len(comparisons),
	})
}
