package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type commentRepository struct {
	db *database.DB
}

func NewCommentRepository(db *database.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, user_id, parent_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		comment.ID, comment.PostID, comment.UserID, comment.ParentID, comment.Content,
		comment.CreatedAt, comment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r *commentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, path, content, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL
	`

	comment := &models.Comment{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID, &comment.Path,
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}

	return comment, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		comment.ID, comment.UserID, comment.Content, comment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or already deleted")
	}

	return nil
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE comments
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or already deleted")
	}

	return nil
}

func (r *commentRepository) ListByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, path, content, created_at, updated_at, deleted_at
		FROM comments
		WHERE post_id = $1 AND deleted_at IS NULL
		ORDER BY path
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments by post ID: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID, &comment.Path,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments: %w", err)
	}

	return comments, nil
}

func (r *commentRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, path, content, created_at, updated_at, deleted_at
		FROM comments
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments by user ID: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID, &comment.Path,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments: %w", err)
	}

	return comments, nil
}

func (r *commentRepository) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, parent_id, path, content, created_at, updated_at, deleted_at
		FROM comments
		WHERE parent_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID, &comment.Path,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate comments: %w", err)
	}

	return comments, nil
}

func (r *commentRepository) CountByPostID(ctx context.Context, postID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM comments
		WHERE post_id = $1 AND deleted_at IS NULL
	`

	var count int
	if err := r.db.GetConnection().QueryRowContext(ctx, query, postID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return count, nil
}
