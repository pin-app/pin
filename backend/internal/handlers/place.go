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

type PlaceHandler struct {
	placeRepo repository.PlaceRepository
	validator *validator.Validate
}

func NewPlaceHandler(placeRepo repository.PlaceRepository) *PlaceHandler {
	return &PlaceHandler{
		placeRepo: placeRepo,
		validator: validator.New(),
	}
}

// CreatePlace handles POST /api/places
func (h *PlaceHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	var req models.PlaceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	place := &models.Place{
		ID:         uuid.New(),
		Name:       req.Name,
		Geometry:   req.Geometry,
		Properties: req.Properties,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.placeRepo.Create(r.Context(), place); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create place"})
		return
	}

	server.WriteJSON(w, http.StatusCreated, place.ToResponse())
}

// GetPlace handles GET /api/places/{id}
func (h *PlaceHandler) GetPlace(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/places/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	place, err := h.placeRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
		return
	}

	server.WriteJSON(w, http.StatusOK, place.ToResponse())
}

// UpdatePlace handles PUT /api/places/{id}
func (h *PlaceHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/places/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	var req models.PlaceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Get existing place
	place, err := h.placeRepo.GetByID(r.Context(), id)
	if err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
		return
	}

	// Update fields
	if req.Name != nil {
		place.Name = *req.Name
	}
	if req.Geometry != nil {
		place.Geometry = req.Geometry
	}
	if req.Properties != nil {
		place.Properties = req.Properties
	}
	place.UpdatedAt = time.Now()

	if err := h.placeRepo.Update(r.Context(), place); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update place"})
		return
	}

	server.WriteJSON(w, http.StatusOK, place.ToResponse())
}

// DeletePlace handles DELETE /api/places/{id}
func (h *PlaceHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/places/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid place ID"})
		return
	}

	if err := h.placeRepo.Delete(r.Context(), id); err != nil {
		server.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Place not found"})
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

// ListPlaces handles GET /api/places
func (h *PlaceHandler) ListPlaces(w http.ResponseWriter, r *http.Request) {
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

	places, err := h.placeRepo.List(r.Context(), limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list places"})
		return
	}

	responses := make([]models.PlaceResponse, len(places))
	for i, place := range places {
		responses[i] = place.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"places": responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

// SearchPlaces handles GET /api/places/search
func (h *PlaceHandler) SearchPlaces(w http.ResponseWriter, r *http.Request) {
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

	places, err := h.placeRepo.Search(r.Context(), query, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to search places"})
		return
	}

	responses := make([]models.PlaceResponse, len(places))
	for i, place := range places {
		responses[i] = place.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"places": responses,
		"query":  query,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

// SearchNearbyPlaces handles GET /api/places/nearby
func (h *PlaceHandler) SearchNearbyPlaces(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	radiusStr := r.URL.Query().Get("radius")

	if latStr == "" || lngStr == "" {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Latitude and longitude parameters are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid latitude"})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid longitude"})
		return
	}

	radius := 10.0 // default 10km
	if radiusStr != "" {
		if r, err := strconv.ParseFloat(radiusStr, 64); err == nil && r > 0 && r <= 100 {
			radius = r
		}
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	places, err := h.placeRepo.SearchNearby(r.Context(), lat, lng, radius, limit)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to search nearby places"})
		return
	}

	responses := make([]models.PlaceResponse, len(places))
	for i, place := range places {
		responses[i] = place.ToResponse()
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"places": responses,
		"lat":    lat,
		"lng":    lng,
		"radius": radius,
		"limit":  limit,
		"count":  len(responses),
	})
}
