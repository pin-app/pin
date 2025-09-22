package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	PlaceID     uuid.UUID  `json:"place_id" db:"place_id"`
	Description *string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type PostImage struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	PostID    uuid.UUID  `json:"post_id" db:"post_id"`
	ImageURL  string     `json:"image_url" db:"image_url"`
	Caption   *string    `json:"caption,omitempty" db:"caption"`
	SortOrder int        `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Comment struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	PostID    uuid.UUID  `json:"post_id" db:"post_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Path      string     `json:"path" db:"path"` // LTREE path for hierarchy
	Content   string     `json:"content" db:"content"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// PostCreateRequest represents the data needed to create a new post
type PostCreateRequest struct {
	PlaceID     uuid.UUID `json:"place_id" validate:"required"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=2000"`
	Images      []string  `json:"images,omitempty" validate:"omitempty,dive,url"`
}

// PostUpdateRequest represents the data that can be updated for a post
type PostUpdateRequest struct {
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
}

// PostResponse represents the post data returned in API responses
type PostResponse struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	PlaceID     uuid.UUID      `json:"place_id"`
	Description *string        `json:"description,omitempty"`
	Images      []PostImage    `json:"images,omitempty"`
	Place       *PlaceResponse `json:"place,omitempty"`
	User        *UserResponse  `json:"user,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ToResponse converts a Post to PostResponse
func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		PlaceID:     p.PlaceID,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// CommentCreateRequest represents the data needed to create a new comment
type CommentCreateRequest struct {
	PostID   uuid.UUID  `json:"post_id" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	Content  string     `json:"content" validate:"required,min=1,max=1000"`
}

// CommentUpdateRequest represents the data that can be updated for a comment
type CommentUpdateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// CommentResponse represents the comment data returned in API responses
type CommentResponse struct {
	ID        uuid.UUID         `json:"id"`
	PostID    uuid.UUID         `json:"post_id"`
	UserID    uuid.UUID         `json:"user_id"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	Content   string            `json:"content"`
	User      *UserResponse     `json:"user,omitempty"`
	Replies   []CommentResponse `json:"replies,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ToResponse converts a Comment to CommentResponse
func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		ParentID:  c.ParentID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// PostImageCreateRequest represents the data needed to create a post image
type PostImageCreateRequest struct {
	PostID    uuid.UUID `json:"post_id" validate:"required"`
	ImageURL  string    `json:"image_url" validate:"required,url"`
	Caption   *string   `json:"caption,omitempty" validate:"omitempty,max=255"`
	SortOrder int       `json:"sort_order" validate:"min=0"`
}
