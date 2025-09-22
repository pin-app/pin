package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type postRepository struct {
	db *database.DB
}

func NewPostRepository(db *database.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (id, user_id, place_id, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		post.ID, post.UserID, post.PlaceID, post.Description, post.CreatedAt, post.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	query := `
		SELECT id, user_id, place_id, description, created_at, updated_at, deleted_at
		FROM posts
		WHERE id = $1 AND deleted_at IS NULL
	`

	post := &models.Post{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.PlaceID, &post.Description,
		&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	return post, nil
}

func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	query := `
		UPDATE posts
		SET description = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		post.ID, post.UserID, post.Description, post.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found or already deleted")
	}

	return nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE posts
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found or already deleted")
	}

	return nil
}

func (r *postRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT id, user_id, place_id, description, created_at, updated_at, deleted_at
		FROM posts
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by user ID: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.UserID, &post.PlaceID, &post.Description,
			&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate posts: %w", err)
	}

	return posts, nil
}

func (r *postRepository) ListByPlaceID(ctx context.Context, placeID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT id, user_id, place_id, description, created_at, updated_at, deleted_at
		FROM posts
		WHERE place_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, placeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by place ID: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.UserID, &post.PlaceID, &post.Description,
			&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate posts: %w", err)
	}

	return posts, nil
}

func (r *postRepository) ListFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	// For now, this is a simple implementation that returns all posts
	// In a real app, this would include posts from followed users, nearby places, etc.
	query := `
		SELECT p.id, p.user_id, p.place_id, p.description, p.created_at, p.updated_at, p.deleted_at
		FROM posts p
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list feed posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID, &post.UserID, &post.PlaceID, &post.Description,
			&post.CreatedAt, &post.UpdatedAt, &post.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate posts: %w", err)
	}

	return posts, nil
}

func (r *postRepository) CreateImage(ctx context.Context, image *models.PostImage) error {
	query := `
		INSERT INTO post_images (id, post_id, image_url, caption, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		image.ID, image.PostID, image.ImageURL, image.Caption, image.SortOrder,
		image.CreatedAt, image.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create post image: %w", err)
	}

	return nil
}

func (r *postRepository) GetImagesByPostID(ctx context.Context, postID uuid.UUID) ([]*models.PostImage, error) {
	query := `
		SELECT id, post_id, image_url, caption, sort_order, created_at, updated_at, deleted_at
		FROM post_images
		WHERE post_id = $1 AND deleted_at IS NULL
		ORDER BY sort_order, created_at
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post images: %w", err)
	}
	defer rows.Close()

	var images []*models.PostImage
	for rows.Next() {
		image := &models.PostImage{}
		err := rows.Scan(
			&image.ID, &image.PostID, &image.ImageURL, &image.Caption, &image.SortOrder,
			&image.CreatedAt, &image.UpdatedAt, &image.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post image: %w", err)
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate post images: %w", err)
	}

	return images, nil
}

func (r *postRepository) UpdateImage(ctx context.Context, image *models.PostImage) error {
	query := `
		UPDATE post_images
		SET caption = $3, sort_order = $4, updated_at = $5
		WHERE id = $1 AND post_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		image.ID, image.PostID, image.Caption, image.SortOrder, image.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update post image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post image not found or already deleted")
	}

	return nil
}

func (r *postRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE post_images
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post image not found or already deleted")
	}

	return nil
}

func (r *postRepository) DeleteImagesByPostID(ctx context.Context, postID uuid.UUID) error {
	query := `
		UPDATE post_images
		SET deleted_at = NOW()
		WHERE post_id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post images: %w", err)
	}

	return nil
}
