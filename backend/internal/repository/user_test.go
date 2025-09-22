package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/models"
)

// InMemoryUserRepository is an in-memory implementation for testing
type InMemoryUserRepository struct {
	users map[string]*models.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *models.User) error {
	r.users[user.ID.String()] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, exists := r.users[id.String()]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range r.users {
		if user.Username != nil && *user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *models.User) error {
	if _, exists := r.users[user.ID.String()]; !exists {
		return ErrUserNotFound
	}
	r.users[user.ID.String()] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := r.users[id.String()]; !exists {
		return ErrUserNotFound
	}
	delete(r.users, id.String())
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	count := 0
	for _, user := range r.users {
		if count >= offset && count < offset+limit {
			users = append(users, user)
		}
		count++
	}
	return users, nil
}

func (r *InMemoryUserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	count := 0
	for _, user := range r.users {
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

func setupTestRepo(t *testing.T) UserRepository {
	return NewInMemoryUserRepository()
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupTestRepo(t)

	user := &models.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Username:    stringPtr("testuser"),
		Bio:         stringPtr("Test bio"),
		Location:    stringPtr("Test City"),
		DisplayName: stringPtr("Test User"),
		PfpURL:      stringPtr("https://example.com/avatar.jpg"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify user was created
	createdUser, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if createdUser.Email != user.Email {
		t.Errorf("GetByID() email = %v, want %v", createdUser.Email, user.Email)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := setupTestRepo(t)

	// Test non-existent user
	id := uuid.New()
	user, err := repo.GetByID(context.Background(), id)
	if err == nil {
		t.Errorf("GetByID() expected error for non-existent user")
	}
	if user != nil {
		t.Errorf("GetByID() expected nil user for non-existent ID")
	}

	// Create a user and test getting it
	user = &models.User{
		ID:        id,
		Email:     "test@example.com",
		Username:  stringPtr("testuser"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Now get the user
	createdUser, err := repo.GetByID(context.Background(), id)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if createdUser == nil {
		t.Errorf("GetByID() expected user")
	}
	if createdUser.Email != user.Email {
		t.Errorf("GetByID() email = %v, want %v", createdUser.Email, user.Email)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := setupTestRepo(t)

	// Test non-existent email
	email := "nonexistent@example.com"
	user, err := repo.GetByEmail(context.Background(), email)
	if err == nil {
		t.Errorf("GetByEmail() expected error for non-existent email")
	}
	if user != nil {
		t.Errorf("GetByEmail() expected nil user for non-existent email")
	}

	// Create a user and test getting it by email
	user = &models.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  stringPtr("testuser"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Now get the user by email
	createdUser, err := repo.GetByEmail(context.Background(), email)
	if err != nil {
		t.Errorf("GetByEmail() error = %v", err)
	}
	if createdUser == nil {
		t.Errorf("GetByEmail() expected user")
	}
	if createdUser.Email != user.Email {
		t.Errorf("GetByEmail() email = %v, want %v", createdUser.Email, user.Email)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	repo := setupTestRepo(t)

	// Test non-existent username
	username := "nonexistentuser"
	user, err := repo.GetByUsername(context.Background(), username)
	if err == nil {
		t.Errorf("GetByUsername() expected error for non-existent username")
	}
	if user != nil {
		t.Errorf("GetByUsername() expected nil user for non-existent username")
	}

	// Create a user and test getting it by username
	user = &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Username:  stringPtr(username),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Now get the user by username
	createdUser, err := repo.GetByUsername(context.Background(), username)
	if err != nil {
		t.Errorf("GetByUsername() error = %v", err)
	}
	if createdUser == nil {
		t.Errorf("GetByUsername() expected user")
	}
	if createdUser.Username == nil || *createdUser.Username != username {
		t.Errorf("GetByUsername() username = %v, want %v", createdUser.Username, username)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := setupTestRepo(t)

	// Test updating non-existent user
	user := &models.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Username:    stringPtr("testuser"),
		Bio:         stringPtr("Updated bio"),
		Location:    stringPtr("Updated City"),
		DisplayName: stringPtr("Updated User"),
		PfpURL:      stringPtr("https://example.com/updated-avatar.jpg"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Update(context.Background(), user)
	if err == nil {
		t.Errorf("Update() expected error for non-existent user")
	}

	// Create a user first
	user.CreatedAt = time.Now()
	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Now update the user
	user.Bio = stringPtr("Really updated bio")
	user.UpdatedAt = time.Now()
	err = repo.Update(context.Background(), user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify the update
	updatedUser, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if updatedUser.Bio == nil || *updatedUser.Bio != "Really updated bio" {
		t.Errorf("Update() bio = %v, want %v", updatedUser.Bio, "Really updated bio")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := setupTestRepo(t)

	// Test deleting non-existent user
	id := uuid.New()
	err := repo.Delete(context.Background(), id)
	if err == nil {
		t.Errorf("Delete() expected error for non-existent user")
	}

	// Create a user first
	user := &models.User{
		ID:        id,
		Email:     "test@example.com",
		Username:  stringPtr("testuser"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Now delete the user
	err = repo.Delete(context.Background(), id)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify the user is deleted (soft delete)
	deletedUser, err := repo.GetByID(context.Background(), id)
	if err == nil {
		t.Errorf("GetByID() expected error for deleted user")
	}
	if deletedUser != nil {
		t.Errorf("GetByID() expected nil user for deleted user")
	}
}

func TestUserRepository_List(t *testing.T) {
	repo := setupTestRepo(t)

	// Create some test users
	users := []*models.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			Username:  stringPtr("user1"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			Username:  stringPtr("user2"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	// Test listing users
	listUsers, err := repo.List(context.Background(), 10, 0)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if len(listUsers) != 2 {
		t.Errorf("List() expected 2 users, got %d", len(listUsers))
	}
}

func TestUserRepository_Search(t *testing.T) {
	repo := setupTestRepo(t)

	// Create some test users
	users := []*models.User{
		{
			ID:        uuid.New(),
			Email:     "testuser1@example.com",
			Username:  stringPtr("testuser1"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "otheruser@example.com",
			Username:  stringPtr("otheruser"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	// Test searching users
	searchUsers, err := repo.Search(context.Background(), "test", 10, 0)
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}
	if len(searchUsers) != 1 {
		t.Errorf("Search() expected 1 user, got %d", len(searchUsers))
	}
	if searchUsers[0].Username == nil || *searchUsers[0].Username != "testuser1" {
		t.Errorf("Search() expected testuser1, got %v", searchUsers[0].Username)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
