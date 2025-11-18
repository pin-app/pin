package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/models"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*models.User, error)
}

// OAuthRepository defines the interface for OAuth-related database operations
type OAuthRepository interface {
	CreateAccount(ctx context.Context, account *models.OAuthAccount) error
	GetAccountByProvider(ctx context.Context, provider models.OAuthProvider, providerID string) (*models.OAuthAccount, error)
	GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.OAuthAccount, error)
	UpdateAccount(ctx context.Context, account *models.OAuthAccount) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error

	CreateState(ctx context.Context, state *models.OAuthState) error
	GetState(ctx context.Context, state string) (*models.OAuthState, error)
	DeleteState(ctx context.Context, state string) error
	CleanupExpiredStates(ctx context.Context) error
}

// SessionRepository defines the interface for session-related database operations
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByToken(ctx context.Context, token string) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	CleanupExpired(ctx context.Context) error
}

// PlaceRepository defines the interface for place-related database operations
type PlaceRepository interface {
	Create(ctx context.Context, place *models.Place) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Place, error)
	Update(ctx context.Context, place *models.Place) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Place, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*models.Place, error)
	SearchNearby(ctx context.Context, lat, lng float64, radiusKm float64, limit int) ([]*models.Place, error)

	CreateRelation(ctx context.Context, relation *models.PlaceRelation) error
	GetRelations(ctx context.Context, placeID uuid.UUID) ([]*models.PlaceRelation, error)
	DeleteRelation(ctx context.Context, id uuid.UUID) error
}

// PostRepository defines the interface for post-related database operations
type PostRepository interface {
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error)
	ListByPlaceID(ctx context.Context, placeID uuid.UUID, limit, offset int) ([]*models.Post, error)
	ListFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error)

	CreateImage(ctx context.Context, image *models.PostImage) error
	GetImagesByPostID(ctx context.Context, postID uuid.UUID) ([]*models.PostImage, error)
	UpdateImage(ctx context.Context, image *models.PostImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
	DeleteImagesByPostID(ctx context.Context, postID uuid.UUID) error
}

// CommentRepository defines the interface for comment-related database operations
type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	CountByPostID(ctx context.Context, postID uuid.UUID) (int, error)
}

// RatingRepository defines the interface for rating-related database operations
type RatingRepository interface {
	CreateRating(ctx context.Context, rating *models.PlaceRating) error
	GetRating(ctx context.Context, userID, placeID uuid.UUID) (*models.PlaceRating, error)
	UpdateRating(ctx context.Context, rating *models.PlaceRating) error
	DeleteRating(ctx context.Context, userID, placeID uuid.UUID) error
	GetRatingsByPlaceID(ctx context.Context, placeID uuid.UUID, limit, offset int) ([]*models.PlaceRating, error)
	GetRatingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.PlaceRating, error)
	GetAverageRating(ctx context.Context, placeID uuid.UUID) (float64, int, error)

	CreateComparison(ctx context.Context, comparison *models.PlaceComparison) error
	GetComparison(ctx context.Context, userID, betterPlaceID, worsePlaceID uuid.UUID) (*models.PlaceComparison, error)
	DeleteComparison(ctx context.Context, userID, betterPlaceID, worsePlaceID uuid.UUID) error
	GetComparisonsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.PlaceComparison, error)
}

// NotificationRepository defines the interface for notification operations
type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	ClearByUserID(ctx context.Context, userID uuid.UUID) error
	SoftDeleteByReference(ctx context.Context, userID, actorID uuid.UUID, notifType models.NotificationType, postID *uuid.UUID, commentID *uuid.UUID) error
}

// LikeRepository defines the interface for post like operations
type LikeRepository interface {
	LikePost(ctx context.Context, postID, userID uuid.UUID) error
	UnlikePost(ctx context.Context, postID, userID uuid.UUID) error
	IsPostLikedByUser(ctx context.Context, postID, userID uuid.UUID) (bool, error)
	CountPostLikes(ctx context.Context, postID uuid.UUID) (int, error)
}
