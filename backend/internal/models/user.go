package models

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderApple  OAuthProvider = "apple"
)

type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Email       string     `json:"email" db:"email"`
	Username    *string    `json:"username,omitempty" db:"username"`
	Bio         *string    `json:"bio,omitempty" db:"bio"`
	Location    *string    `json:"location,omitempty" db:"location"`
	DisplayName *string    `json:"display_name,omitempty" db:"display_name"`
	PfpURL      *string    `json:"pfp_url,omitempty" db:"pfp_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type OAuthAccount struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	UserID         uuid.UUID     `json:"user_id" db:"user_id"`
	Provider       OAuthProvider `json:"provider" db:"provider"`
	ProviderID     string        `json:"provider_id" db:"provider_id"`
	ProviderEmail  *string       `json:"provider_email,omitempty" db:"provider_email"`
	ProviderName   *string       `json:"provider_name,omitempty" db:"provider_name"`
	AccessToken    *string       `json:"access_token,omitempty" db:"access_token"`
	RefreshToken   *string       `json:"refresh_token,omitempty" db:"refresh_token"`
	TokenExpiresAt *time.Time    `json:"token_expires_at,omitempty" db:"token_expires_at"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Session struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	SessionToken string     `json:"session_token" db:"session_token"`
	ExpiresAt    time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type OAuthState struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	State        string        `json:"state" db:"state"`
	CodeVerifier *string       `json:"code_verifier,omitempty" db:"code_verifier"`
	Provider     OAuthProvider `json:"provider" db:"provider"`
	RedirectURL  *string       `json:"redirect_url,omitempty" db:"redirect_url"`
	ExpiresAt    time.Time     `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
}

// UserCreateRequest represents the data needed to create a new user
type UserCreateRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Username    *string `json:"username,omitempty" validate:"omitempty,min=3,max=30,alphanum"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=100"`
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,max=100"`
	PfpURL      *string `json:"pfp_url,omitempty" validate:"omitempty,url"`
}

// UserUpdateRequest represents the data that can be updated for a user
type UserUpdateRequest struct {
	Username    *string `json:"username,omitempty" validate:"omitempty,min=3,max=30,alphanum"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=100"`
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,max=100"`
	PfpURL      *string `json:"pfp_url,omitempty" validate:"omitempty,url"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Username    *string   `json:"username,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	Location    *string   `json:"location,omitempty"`
	DisplayName *string   `json:"display_name,omitempty"`
	PfpURL      *string   `json:"pfp_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		Bio:         u.Bio,
		Location:    u.Location,
		DisplayName: u.DisplayName,
		PfpURL:      u.PfpURL,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
