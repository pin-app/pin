package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/middleware"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/server"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	oauthRepo    repository.OAuthRepository
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	authMW       *middleware.AuthMiddleware
	googleConfig *oauth2.Config
	appleConfig  *oauth2.Config
}

type OAuthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type OAuthResponse struct {
	SessionToken string              `json:"session_token"`
	User         models.UserResponse `json:"user"`
	ExpiresAt    time.Time           `json:"expires_at"`
}

func NewOAuthHandler(oauthRepo repository.OAuthRepository, userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *OAuthHandler {
	authMW := middleware.NewAuthMiddleware(sessionRepo, userRepo)

	// Google OAuth configuration
	googleConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Apple OAuth configuration
	appleConfig := &oauth2.Config{
		ClientID:     os.Getenv("APPLE_CLIENT_ID"),
		ClientSecret: os.Getenv("APPLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("APPLE_REDIRECT_URL"),
		Scopes:       []string{"name", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
	}

	return &OAuthHandler{
		oauthRepo:    oauthRepo,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		authMW:       authMW,
		googleConfig: googleConfig,
		appleConfig:  appleConfig,
	}
}

// Google OAuth

func (h *OAuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := h.generateState(models.OAuthProviderGoogle, r.URL.Query().Get("redirect_url"))
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate state"})
		return
	}

	authURL := h.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing code or state parameter"})
		return
	}

	// Verify state
	oauthState, err := h.oauthRepo.GetState(r.Context(), state)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid or expired state"})
		return
	}

	// Exchange code for token
	token, err := h.googleConfig.Exchange(r.Context(), code)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to exchange code for token"})
		return
	}

	// Get user info from Google
	userInfo, err := h.getGoogleUserInfo(r.Context(), token.AccessToken)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get user info"})
		return
	}

	// Process OAuth account
	redirectURL := ""
	if oauthState.RedirectURL != nil {
		redirectURL = *oauthState.RedirectURL
	}
	response, err := h.processOAuthAccount(r.Context(), models.OAuthProviderGoogle, userInfo, token, redirectURL)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Clean up state
	_ = h.oauthRepo.DeleteState(r.Context(), state)

	server.WriteJSON(w, http.StatusOK, response)
}

// Apple OAuth

func (h *OAuthHandler) AppleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := h.generateState(models.OAuthProviderApple, r.URL.Query().Get("redirect_url"))
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate state"})
		return
	}

	authURL := h.appleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) AppleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing code or state parameter"})
		return
	}

	// Verify state
	oauthState, err := h.oauthRepo.GetState(r.Context(), state)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid or expired state"})
		return
	}

	// Exchange code for token
	token, err := h.appleConfig.Exchange(r.Context(), code)
	if err != nil {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to exchange code for token"})
		return
	}

	// Get user info from Apple
	userInfo, err := h.getAppleUserInfo(r.Context(), token.AccessToken)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get user info"})
		return
	}

	// Process OAuth account
	redirectURL := ""
	if oauthState.RedirectURL != nil {
		redirectURL = *oauthState.RedirectURL
	}
	response, err := h.processOAuthAccount(r.Context(), models.OAuthProviderApple, userInfo, token, redirectURL)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Clean up state
	_ = h.oauthRepo.DeleteState(r.Context(), state)

	server.WriteJSON(w, http.StatusOK, response)
}

// Logout

func (h *OAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Authorization header required"})
		return
	}

	parts := []string{}
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		parts = []string{"Bearer", authHeader[7:]}
	} else {
		parts = []string{"Bearer", authHeader}
	}

	if len(parts) != 2 || parts[0] != "Bearer" {
		server.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid authorization header format"})
		return
	}

	sessionToken := parts[1]
	if err := h.authMW.DeleteSession(r.Context(), sessionToken); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
		return
	}

	server.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// Helper methods

func (h *OAuthHandler) generateState(provider models.OAuthProvider, redirectURL string) (string, error) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", err
	}
	state := hex.EncodeToString(stateBytes)

	oauthState := &models.OAuthState{
		ID:          uuid.New(),
		State:       state,
		Provider:    provider,
		RedirectURL: &redirectURL,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		CreatedAt:   time.Now(),
	}

	if err := h.oauthRepo.CreateState(context.Background(), oauthState); err != nil {
		return "", err
	}

	return state, nil
}

func (h *OAuthHandler) getGoogleUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	client := h.googleConfig.Client(ctx, &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (h *OAuthHandler) getAppleUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	// Apple doesn't provide a user info endpoint, so we need to decode the ID token
	// For simplicity, we'll create a basic user info structure
	// In a real implementation, you'd decode the JWT ID token
	return map[string]interface{}{
		"id":    "apple_user_" + accessToken[:8], // Simplified
		"email": "user@privaterelay.appleid.com", // Apple uses private relay
		"name":  "Apple User",
	}, nil
}

func (h *OAuthHandler) processOAuthAccount(ctx context.Context, provider models.OAuthProvider, userInfo map[string]interface{}, token *oauth2.Token, redirectURL string) (*OAuthResponse, error) {
	providerID := h.getStringFromMap(userInfo, "id")
	email := h.getStringFromMap(userInfo, "email")
	name := h.getStringFromMap(userInfo, "name")

	if providerID == "" {
		return nil, fmt.Errorf("provider ID not found in user info")
	}

	// Check if OAuth account already exists
	existingAccount, err := h.oauthRepo.GetAccountByProvider(ctx, provider, providerID)
	if err == nil {
		// Update existing account
		existingAccount.AccessToken = &token.AccessToken
		if token.RefreshToken != "" {
			existingAccount.RefreshToken = &token.RefreshToken
		}
		if !token.Expiry.IsZero() {
			existingAccount.TokenExpiresAt = &token.Expiry
		}
		existingAccount.UpdatedAt = time.Now()

		if err := h.oauthRepo.UpdateAccount(ctx, existingAccount); err != nil {
			return nil, fmt.Errorf("failed to update OAuth account: %w", err)
		}

		// Get user
		user, err := h.userRepo.GetByID(ctx, existingAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// Create session
		session, err := h.authMW.CreateSession(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}

		return &OAuthResponse{
			SessionToken: session.SessionToken,
			User:         user.ToResponse(),
			ExpiresAt:    session.ExpiresAt,
		}, nil
	}

	// Create new user and OAuth account
	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if name != "" {
		user.DisplayName = &name
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create OAuth account
	oauthAccount := &models.OAuthAccount{
		ID:            uuid.New(),
		UserID:        user.ID,
		Provider:      provider,
		ProviderID:    providerID,
		ProviderEmail: &email,
		ProviderName:  &name,
		AccessToken:   &token.AccessToken,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if token.RefreshToken != "" {
		oauthAccount.RefreshToken = &token.RefreshToken
	}
	if !token.Expiry.IsZero() {
		oauthAccount.TokenExpiresAt = &token.Expiry
	}

	if err := h.oauthRepo.CreateAccount(ctx, oauthAccount); err != nil {
		return nil, fmt.Errorf("failed to create OAuth account: %w", err)
	}

	// Create session
	session, err := h.authMW.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &OAuthResponse{
		SessionToken: session.SessionToken,
		User:         user.ToResponse(),
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

func (h *OAuthHandler) getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
