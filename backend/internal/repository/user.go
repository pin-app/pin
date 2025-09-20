package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, username, bio, location, display_name, pfp_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		user.ID, user.Email, user.Username, user.Bio, user.Location,
		user.DisplayName, user.PfpURL, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, username, bio, location, display_name, pfp_url, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.Bio, &user.Location,
		&user.DisplayName, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, bio, location, display_name, pfp_url, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.Bio, &user.Location,
		&user.DisplayName, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, email, username, bio, location, display_name, pfp_url, created_at, updated_at, deleted_at
		FROM users
		WHERE LOWER(username) = LOWER($1) AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.Bio, &user.Location,
		&user.DisplayName, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $2, username = $3, bio = $4, location = $5, display_name = $6, pfp_url = $7, updated_at = $8
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		user.ID, user.Email, user.Username, user.Bio, user.Location,
		user.DisplayName, user.PfpURL, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, email, username, bio, location, display_name, pfp_url, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.Bio, &user.Location,
			&user.DisplayName, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}

func (r *userRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.User, error) {
	searchQuery := `
		SELECT id, email, username, bio, location, display_name, pfp_url, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		AND (
			LOWER(username) LIKE LOWER($1) OR
			LOWER(display_name) LIKE LOWER($1) OR
			LOWER(bio) LIKE LOWER($1)
		)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.GetConnection().QueryContext(ctx, searchQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.Bio, &user.Location,
			&user.DisplayName, &user.PfpURL, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}
