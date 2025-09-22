package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type FollowRepository interface {
	CreateFollow(ctx context.Context, follow *models.Follow) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	GetFollow(ctx context.Context, followerID, followingID uuid.UUID) (*models.Follow, error)
	ListFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.User, error)
	ListFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.User, error)
	GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetFollowersCount(ctx context.Context, userID uuid.UUID) (int, error)
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (*models.UserStats, error)
}

type followRepository struct {
	db *database.DB
}

func NewFollowRepository(db *database.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) CreateFollow(ctx context.Context, follow *models.Follow) error {
	query := `
		INSERT INTO follows (id, follower_id, following_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.GetConnection().ExecContext(ctx, query, follow.ID, follow.FollowerID, follow.FollowingID, follow.CreatedAt, follow.UpdatedAt)
	return err
}

func (r *followRepository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	result, err := r.db.GetConnection().ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follow relationship not found")
	}

	return nil
}

func (r *followRepository) GetFollow(ctx context.Context, followerID, followingID uuid.UUID) (*models.Follow, error) {
	query := `
		SELECT id, follower_id, following_id, created_at, updated_at
		FROM follows
		WHERE follower_id = $1 AND following_id = $2
	`

	var follow models.Follow
	err := r.db.GetConnection().QueryRowContext(ctx, query, followerID, followingID).Scan(
		&follow.ID,
		&follow.FollowerID,
		&follow.FollowingID,
		&follow.CreatedAt,
		&follow.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("follow relationship not found")
		}
		return nil, err
	}

	return &follow, nil
}

func (r *followRepository) ListFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.bio, u.location, u.display_name, u.pfp_url, u.created_at, u.updated_at
		FROM follows f
		JOIN users u ON f.following_id = u.id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Bio,
			&user.Location,
			&user.DisplayName,
			&user.PfpURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *followRepository) ListFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.bio, u.location, u.display_name, u.pfp_url, u.created_at, u.updated_at
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Bio,
			&user.Location,
			&user.DisplayName,
			&user.PfpURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *followRepository) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`
	var count int
	err := r.db.GetConnection().QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *followRepository) GetFollowersCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM follows WHERE following_id = $1`
	var count int
	err := r.db.GetConnection().QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`
	var exists bool
	err := r.db.GetConnection().QueryRowContext(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *followRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*models.UserStats, error) {
	// Get posts count
	postsCountQuery := `SELECT COUNT(*) FROM posts WHERE user_id = $1`
	var postsCount int
	err := r.db.GetConnection().QueryRowContext(ctx, postsCountQuery, userID).Scan(&postsCount)
	if err != nil {
		return nil, err
	}

	// Get following count
	followingCount, err := r.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get followers count
	followersCount, err := r.GetFollowersCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserStats{
		UserID:         userID,
		PostsCount:     postsCount,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
	}, nil
}
