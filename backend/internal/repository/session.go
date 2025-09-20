package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type sessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, session_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		session.ID, session.UserID, session.SessionToken, session.ExpiresAt,
		session.CreatedAt, session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, session_token, expires_at, created_at, updated_at, deleted_at
		FROM sessions
		WHERE session_token = $1 AND expires_at > NOW() AND deleted_at IS NULL
	`

	session := &models.Session{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.SessionToken, &session.ExpiresAt,
		&session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session by token: %w", err)
	}

	return session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	query := `
		SELECT id, user_id, session_token, expires_at, created_at, updated_at, deleted_at
		FROM sessions
		WHERE user_id = $1 AND expires_at > NOW() AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by user ID: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.UserID, &session.SessionToken, &session.ExpiresAt,
			&session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return sessions, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE sessions
		SET expires_at = $3, updated_at = $4
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		session.ID, session.UserID, session.ExpiresAt, session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already deleted")
	}

	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE sessions
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already deleted")
	}

	return nil
}

func (r *sessionRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `
		UPDATE sessions
		SET deleted_at = NOW()
		WHERE session_token = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session by token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already deleted")
	}

	return nil
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE sessions
		SET deleted_at = NOW()
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user ID: %w", err)
	}

	return nil
}

func (r *sessionRepository) CleanupExpired(ctx context.Context) error {
	query := `
		UPDATE sessions
		SET deleted_at = NOW()
		WHERE expires_at <= NOW() AND deleted_at IS NULL
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return nil
}
