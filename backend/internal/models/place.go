package models

import (
	"time"

	"github.com/google/uuid"
)

type PlaceRelationType string

const (
	PlaceRelationContains PlaceRelationType = "CONTAINS"
	PlaceRelationPartOf   PlaceRelationType = "PART_OF"
	PlaceRelationOverlaps PlaceRelationType = "OVERLAPS"
)

type Place struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	Name       string         `json:"name" db:"name"`
	Geometry   *string        `json:"geometry,omitempty" db:"geometry"` // PostGIS geometry as WKT
	Properties map[string]any `json:"properties,omitempty" db:"properties"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
}

type PlaceRelation struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	FromPlaceID  uuid.UUID         `json:"from_place_id" db:"from_place_id"`
	ToPlaceID    uuid.UUID         `json:"to_place_id" db:"to_place_id"`
	RelationType PlaceRelationType `json:"relation_type" db:"relation_type"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`
}

type PlaceRating struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	PlaceID   uuid.UUID  `json:"place_id" db:"place_id"`
	Rating    int        `json:"rating" db:"rating"` // 0-100
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type PlaceComparison struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	BetterPlaceID uuid.UUID  `json:"better_place_id" db:"better_place_id"`
	WorsePlaceID  uuid.UUID  `json:"worse_place_id" db:"worse_place_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// PlaceCreateRequest represents the data needed to create a new place
type PlaceCreateRequest struct {
	Name       string         `json:"name" validate:"required,min=1,max=255"`
	Geometry   *string        `json:"geometry,omitempty" validate:"omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// PlaceUpdateRequest represents the data that can be updated for a place
type PlaceUpdateRequest struct {
	Name       *string        `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Geometry   *string        `json:"geometry,omitempty" validate:"omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// PlaceResponse represents the place data returned in API responses
type PlaceResponse struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Geometry   *string        `json:"geometry,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// ToResponse converts a Place to PlaceResponse
func (p *Place) ToResponse() PlaceResponse {
	return PlaceResponse{
		ID:         p.ID,
		Name:       p.Name,
		Geometry:   p.Geometry,
		Properties: p.Properties,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// PlaceRatingRequest represents the data needed to create/update a place rating
type PlaceRatingRequest struct {
	Rating int `json:"rating" validate:"required,min=0,max=100"`
}

// PlaceComparisonRequest represents the data needed to create a place comparison
type PlaceComparisonRequest struct {
	BetterPlaceID uuid.UUID `json:"better_place_id" validate:"required"`
	WorsePlaceID  uuid.UUID `json:"worse_place_id" validate:"required"`
}
