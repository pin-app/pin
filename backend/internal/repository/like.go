package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
)

type likeRepository struct {
	db *database.DB
}

func NewLikeRepository(db *database.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) LikePost(ctx context.Context, postID, userID uuid.UUID) error {
	query := `
		INSERT INTO post_likes (id, post_id, user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, user_id) DO NOTHING
	`

	if _, err := r.db.GetConnection().ExecContext(ctx, query, uuid.New(), postID, userID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

func (r *likeRepository) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	query := `
		DELETE FROM post_likes
		WHERE post_id = $1 AND user_id = $2
	`

	if _, err := r.db.GetConnection().ExecContext(ctx, query, postID, userID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	return nil
}

func (r *likeRepository) IsPostLikedByUser(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM post_likes
			WHERE post_id = $1 AND user_id = $2
		)
	`

	var exists bool
	if err := r.db.GetConnection().QueryRowContext(ctx, query, postID, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}

	return exists, nil
}

func (r *likeRepository) CountPostLikes(ctx context.Context, postID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM post_likes
		WHERE post_id = $1
	`

	var count int
	if err := r.db.GetConnection().QueryRowContext(ctx, query, postID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count likes: %w", err)
	}

	return count, nil
}
