package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/models"
)

type oauthRepository struct {
	db *database.DB
}

func NewOAuthRepository(db *database.DB) OAuthRepository {
	return &oauthRepository{db: db}
}

func (r *oauthRepository) CreateAccount(ctx context.Context, account *models.OAuthAccount) error {
	query := `
		INSERT INTO oauth_accounts (id, user_id, provider, provider_id, provider_email, provider_name, 
			access_token, refresh_token, token_expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		account.ID, account.UserID, account.Provider, account.ProviderID, account.ProviderEmail,
		account.ProviderName, account.AccessToken, account.RefreshToken, account.TokenExpiresAt,
		account.CreatedAt, account.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create OAuth account: %w", err)
	}

	return nil
}

func (r *oauthRepository) GetAccountByProvider(ctx context.Context, provider models.OAuthProvider, providerID string) (*models.OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_id, provider_email, provider_name,
			access_token, refresh_token, token_expires_at, created_at, updated_at, deleted_at
		FROM oauth_accounts
		WHERE provider = $1 AND provider_id = $2 AND deleted_at IS NULL
	`

	account := &models.OAuthAccount{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, provider, providerID).Scan(
		&account.ID, &account.UserID, &account.Provider, &account.ProviderID, &account.ProviderEmail,
		&account.ProviderName, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt,
		&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("OAuth account not found")
		}
		return nil, fmt.Errorf("failed to get OAuth account: %w", err)
	}

	return account, nil
}

func (r *oauthRepository) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_id, provider_email, provider_name,
			access_token, refresh_token, token_expires_at, created_at, updated_at, deleted_at
		FROM oauth_accounts
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.GetConnection().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*models.OAuthAccount
	for rows.Next() {
		account := &models.OAuthAccount{}
		err := rows.Scan(
			&account.ID, &account.UserID, &account.Provider, &account.ProviderID, &account.ProviderEmail,
			&account.ProviderName, &account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt,
			&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan OAuth account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate OAuth accounts: %w", err)
	}

	return accounts, nil
}

func (r *oauthRepository) UpdateAccount(ctx context.Context, account *models.OAuthAccount) error {
	query := `
		UPDATE oauth_accounts
		SET provider_email = $3, provider_name = $4, access_token = $5, refresh_token = $6,
			token_expires_at = $7, updated_at = $8
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query,
		account.ID, account.UserID, account.ProviderEmail, account.ProviderName,
		account.AccessToken, account.RefreshToken, account.TokenExpiresAt, account.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update OAuth account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("OAuth account not found or already deleted")
	}

	return nil
}

func (r *oauthRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE oauth_accounts
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.GetConnection().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete OAuth account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("OAuth account not found or already deleted")
	}

	return nil
}

func (r *oauthRepository) CreateState(ctx context.Context, state *models.OAuthState) error {
	query := `
		INSERT INTO oauth_states (id, state, code_verifier, provider, redirect_url, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.GetConnection().ExecContext(ctx, query,
		state.ID, state.State, state.CodeVerifier, state.Provider, state.RedirectURL,
		state.ExpiresAt, state.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create OAuth state: %w", err)
	}

	return nil
}

func (r *oauthRepository) GetState(ctx context.Context, state string) (*models.OAuthState, error) {
	query := `
		SELECT id, state, code_verifier, provider, redirect_url, expires_at, created_at
		FROM oauth_states
		WHERE state = $1 AND expires_at > NOW()
	`

	oauthState := &models.OAuthState{}
	err := r.db.GetConnection().QueryRowContext(ctx, query, state).Scan(
		&oauthState.ID, &oauthState.State, &oauthState.CodeVerifier, &oauthState.Provider,
		&oauthState.RedirectURL, &oauthState.ExpiresAt, &oauthState.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("OAuth state not found or expired")
		}
		return nil, fmt.Errorf("failed to get OAuth state: %w", err)
	}

	return oauthState, nil
}

func (r *oauthRepository) DeleteState(ctx context.Context, state string) error {
	query := `DELETE FROM oauth_states WHERE state = $1`

	result, err := r.db.GetConnection().ExecContext(ctx, query, state)
	if err != nil {
		return fmt.Errorf("failed to delete OAuth state: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("OAuth state not found")
	}

	return nil
}

func (r *oauthRepository) CleanupExpiredStates(ctx context.Context) error {
	query := `DELETE FROM oauth_states WHERE expires_at <= NOW()`

	_, err := r.db.GetConnection().ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired OAuth states: %w", err)
	}

	return nil
}
