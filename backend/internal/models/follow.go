package models

import (
	"time"

	"github.com/google/uuid"
)

// Follow represents a follow relationship between users
type Follow struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FollowerID  uuid.UUID `json:"follower_id" db:"follower_id"`   // User who is following
	FollowingID uuid.UUID `json:"following_id" db:"following_id"` // User being followed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FollowRequest represents the data needed to create a follow relationship
type FollowRequest struct {
	FollowingID uuid.UUID `json:"following_id" validate:"required"`
}

// FollowResponse represents a follow relationship in API responses
type FollowResponse struct {
	ID          uuid.UUID `json:"id"`
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserStats represents user statistics
type UserStats struct {
	UserID         uuid.UUID `json:"user_id"`
	PostsCount     int       `json:"posts_count"`
	FollowingCount int       `json:"following_count"`
	FollowersCount int       `json:"followers_count"`
}

// ToResponse converts a Follow to FollowResponse
func (f *Follow) ToResponse() FollowResponse {
	return FollowResponse{
		ID:          f.ID,
		FollowerID:  f.FollowerID,
		FollowingID: f.FollowingID,
		CreatedAt:   f.CreatedAt,
	}
}
