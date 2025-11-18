package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeLikePost     NotificationType = "like_post"
	NotificationTypeCommentPost  NotificationType = "comment_post"
	NotificationTypeCommentReply NotificationType = "comment_reply"
)

type Notification struct {
	ID        uuid.UUID         `json:"id" db:"id"`
	UserID    uuid.UUID         `json:"user_id" db:"user_id"`
	ActorID   uuid.UUID         `json:"actor_id" db:"actor_id"`
	PostID    *uuid.UUID        `json:"post_id,omitempty" db:"post_id"`
	CommentID *uuid.UUID        `json:"comment_id,omitempty" db:"comment_id"`
	Type      NotificationType  `json:"type" db:"type"`
	Data      map[string]string `json:"data" db:"data"`
	ReadAt    *time.Time        `json:"read_at,omitempty" db:"read_at"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`
}

type NotificationResponse struct {
	ID        uuid.UUID         `json:"id"`
	Type      NotificationType  `json:"type"`
	PostID    *uuid.UUID        `json:"post_id,omitempty"`
	CommentID *uuid.UUID        `json:"comment_id,omitempty"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
	ReadAt    *time.Time        `json:"read_at,omitempty"`
	Actor     *UserResponse     `json:"actor,omitempty"`
}

func (n *Notification) ToResponse() NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		Type:      n.Type,
		PostID:    n.PostID,
		CommentID: n.CommentID,
		Data:      n.Data,
		CreatedAt: n.CreatedAt,
		ReadAt:    n.ReadAt,
	}
}
