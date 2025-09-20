package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
)

type AuthMiddleware struct {
	sessionRepo repository.SessionRepository
	userRepo    repository.UserRepository
	devMode     bool
}

type AuthContextKey string

const (
	UserIDKey    AuthContextKey = "user_id"
	SessionKey   AuthContextKey = "session"
	IsDevModeKey AuthContextKey = "is_dev_mode"
)

func NewAuthMiddleware(sessionRepo repository.SessionRepository, userRepo repository.UserRepository) *AuthMiddleware {
	devMode := false
	if devModeStr := os.Getenv("DEV_MODE"); devModeStr != "" {
		if parsed, err := strconv.ParseBool(devModeStr); err == nil {
			devMode = parsed
		}
	}

	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		devMode:     devMode,
	}
}

func (a *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), IsDevModeKey, a.devMode)
		r = r.WithContext(ctx)

		// In dev mode, create a mock user if no auth is provided
		if a.devMode {
			userID := a.getDevUserID(r)
			if userID != uuid.Nil {
				ctx = context.WithValue(ctx, UserIDKey, userID)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		// Extract session token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		// Check for Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		sessionToken := parts[1]
		if sessionToken == "" {
			http.Error(w, `{"error": "Session token required"}`, http.StatusUnauthorized)
			return
		}

		// Get session from database
		session, err := a.sessionRepo.GetByToken(r.Context(), sessionToken)
		if err != nil {
			http.Error(w, `{"error": "Invalid or expired session"}`, http.StatusUnauthorized)
			return
		}

		// Add user ID and session to context
		ctx = context.WithValue(ctx, UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, SessionKey, session)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), IsDevModeKey, a.devMode)
		r = r.WithContext(ctx)

		// In dev mode, try to get dev user ID
		if a.devMode {
			userID := a.getDevUserID(r)
			if userID != uuid.Nil {
				ctx = context.WithValue(ctx, UserIDKey, userID)
				r = r.WithContext(ctx)
			}
		} else {
			// Try to extract session token
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" && parts[1] != "" {
					session, err := a.sessionRepo.GetByToken(r.Context(), parts[1])
					if err == nil {
						ctx = context.WithValue(ctx, UserIDKey, session.UserID)
						ctx = context.WithValue(ctx, SessionKey, session)
					}
				}
			}
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (a *AuthMiddleware) getDevUserID(r *http.Request) uuid.UUID {
	// In dev mode, check for a dev-user-id header or query param
	devUserID := r.Header.Get("X-Dev-User-ID")
	if devUserID == "" {
		devUserID = r.URL.Query().Get("dev_user_id")
	}

	if devUserID != "" {
		if userID, err := uuid.Parse(devUserID); err == nil {
			return userID
		}
	}

	// If no dev user ID provided, try to get or create a default dev user
	return a.getOrCreateDevUser(r.Context())
}

func (a *AuthMiddleware) getOrCreateDevUser(ctx context.Context) uuid.UUID {
	// Try to find an existing dev user
	devUser, err := a.userRepo.GetByEmail(ctx, "dev@localhost")
	if err == nil {
		return devUser.ID
	}

	// Create a new dev user
	devUser = &models.User{
		ID:          uuid.New(),
		Email:       "dev@localhost",
		Username:    stringPtr("devuser"),
		DisplayName: stringPtr("Dev User"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := a.userRepo.Create(ctx, devUser); err != nil {
		// If we can't create a dev user, return nil UUID
		return uuid.Nil
	}

	return devUser.ID
}

func (a *AuthMiddleware) CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error) {
	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	session := &models.Session{
		ID:           uuid.New(),
		UserID:       userID,
		SessionToken: sessionToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := a.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (a *AuthMiddleware) DeleteSession(ctx context.Context, sessionToken string) error {
	return a.sessionRepo.DeleteByToken(ctx, sessionToken)
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func stringPtr(s string) *string {
	return &s
}

// Helper functions to extract data from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func GetSessionFromContext(ctx context.Context) (*models.Session, bool) {
	session, ok := ctx.Value(SessionKey).(*models.Session)
	return session, ok
}

func IsDevModeFromContext(ctx context.Context) bool {
	isDevMode, ok := ctx.Value(IsDevModeKey).(bool)
	return ok && isDevMode
}
