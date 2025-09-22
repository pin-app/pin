package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[uuid.UUID]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uuid.UUID]*models.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username != nil && *user.Username == username {
			return user, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return repository.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.users[id]; !exists {
		return repository.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	count := 0
	for _, user := range m.users {
		if count >= offset && count < offset+limit {
			users = append(users, user)
		}
		count++
	}
	return users, nil
}

func (m *MockUserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	count := 0
	for _, user := range m.users {
		if count >= offset && count < offset+limit {
			// Simple search implementation
			if user.Username != nil && contains(*user.Username, query) {
				users = append(users, user)
			}
		}
		count++
	}
	return users, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	reqBody := models.UserCreateRequest{
		Email:       "test@example.com",
		Username:    stringPtr("testuser"),
		Bio:         stringPtr("Test bio"),
		Location:    stringPtr("Test City"),
		DisplayName: stringPtr("Test User"),
		PfpURL:      stringPtr("https://example.com/avatar.jpg"),
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("CreateUser() status = %v, want %v", rr.Code, http.StatusCreated)
	}

	var response models.UserResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("CreateUser() failed to decode response: %v", err)
	}

	if response.Email != reqBody.Email {
		t.Errorf("CreateUser() email = %v, want %v", response.Email, reqBody.Email)
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	// Create a test user
	userID := uuid.New()
	user := &models.User{
		ID:          userID,
		Email:       "test@example.com",
		Username:    stringPtr("testuser"),
		Bio:         stringPtr("Test bio"),
		Location:    stringPtr("Test City"),
		DisplayName: stringPtr("Test User"),
		PfpURL:      stringPtr("https://example.com/avatar.jpg"),
	}
	mockRepo.Create(context.Background(), user)

	req := httptest.NewRequest("GET", "/api/users/"+userID.String(), nil)
	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GetUser() status = %v, want %v", rr.Code, http.StatusOK)
	}

	var response models.UserResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("GetUser() failed to decode response: %v", err)
	}

	if response.ID != userID {
		t.Errorf("GetUser() ID = %v, want %v", response.ID, userID)
	}
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	nonExistentID := uuid.New()
	req := httptest.NewRequest("GET", "/api/users/"+nonExistentID.String(), nil)
	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GetUser() status = %v, want %v", rr.Code, http.StatusNotFound)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	// Create a test user
	userID := uuid.New()
	user := &models.User{
		ID:          userID,
		Email:       "test@example.com",
		Username:    stringPtr("testuser"),
		Bio:         stringPtr("Test bio"),
		Location:    stringPtr("Test City"),
		DisplayName: stringPtr("Test User"),
		PfpURL:      stringPtr("https://example.com/avatar.jpg"),
	}
	mockRepo.Create(context.Background(), user)

	reqBody := models.UserUpdateRequest{
		Bio:      stringPtr("Updated bio"),
		Location: stringPtr("Updated City"),
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/users/"+userID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("UpdateUser() status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	// Create a test user
	userID := uuid.New()
	user := &models.User{
		ID:    userID,
		Email: "test@example.com",
	}
	mockRepo.Create(context.Background(), user)

	req := httptest.NewRequest("DELETE", "/api/users/"+userID.String(), nil)
	rr := httptest.NewRecorder()
	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("DeleteUser() status = %v, want %v", rr.Code, http.StatusNoContent)
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	req := httptest.NewRequest("GET", "/api/users", nil)
	rr := httptest.NewRecorder()
	handler.ListUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ListUsers() status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestUserHandler_SearchUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	handler := NewUserHandler(mockRepo)

	req := httptest.NewRequest("GET", "/api/users/search?q=test", nil)
	rr := httptest.NewRecorder()
	handler.SearchUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("SearchUsers() status = %v, want %v", rr.Code, http.StatusOK)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
