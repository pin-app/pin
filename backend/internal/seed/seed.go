package seed

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
)

type PostImage struct {
	ImageURL  string
	Caption   string
	SortOrder int
}

type Comment struct {
	Username string
	Content  string
	Replies  []Comment
}

type Post struct {
	Description string
	PlaceName   string // Name of the place where this post was made
	Images      []PostImage
	Comments    []Comment
	CreatedAt   time.Time
}

type UserProfile struct {
	Email       string
	Username    string
	DisplayName string
	Bio         string
	Location    string
	PfpURL      string
	Posts       []Post
	CreatedAt   time.Time
}

type PlaceProfile struct {
	Name       string
	Geometry   string
	Properties map[string]any
}

type FollowRelationship struct {
	FollowerUsername  string
	FollowingUsername string
}

type Seeder struct {
	userRepo    repository.UserRepository
	placeRepo   repository.PlaceRepository
	postRepo    repository.PostRepository
	ratingRepo  repository.RatingRepository
	commentRepo repository.CommentRepository
	followRepo  repository.FollowRepository
}

func NewSeeder(
	userRepo repository.UserRepository,
	placeRepo repository.PlaceRepository,
	postRepo repository.PostRepository,
	ratingRepo repository.RatingRepository,
	commentRepo repository.CommentRepository,
	followRepo repository.FollowRepository,
) *Seeder {
	return &Seeder{
		userRepo:    userRepo,
		placeRepo:   placeRepo,
		postRepo:    postRepo,
		ratingRepo:  ratingRepo,
		commentRepo: commentRepo,
		followRepo:  followRepo,
	}
}

func (s *Seeder) SeedDevData(ctx context.Context) error {
	slog.Info("seeding development data")

	userProfiles := s.getUserProfiles()
	placeProfiles := s.getPlaceProfiles()
	followRelationships := s.getFollowRelationships()

	places, err := s.seedPlaces(ctx, placeProfiles)
	if err != nil {
		return fmt.Errorf("failed to seed places: %w", err)
	}

	users, err := s.seedUserProfiles(ctx, userProfiles, places)
	if err != nil {
		return fmt.Errorf("failed to seed user profiles: %w", err)
	}

	if err := s.seedFollowRelationships(ctx, users, followRelationships); err != nil {
		return fmt.Errorf("failed to seed follows: %w", err)
	}

	slog.Info("development data seeding complete",
		"users", len(users),
		"places", len(places),
	)

	return nil
}

// data definition functions - modify these to add/change user profiles
func (s *Seeder) getUserProfiles() []UserProfile {
	return []UserProfile{
		{
			Email:       "acc@raquent.in",
			Username:    "raquentin",
			DisplayName: "Race",
			Bio:         "6'2 feminist, matcha labubu keshi rave",
			Location:    "uhouse",
			PfpURL:      "https://i.imgur.com/ZpA6U3C.jpeg",
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			Posts: []Post{
				{
					Description: "Ts calm, good place to live",
					PlaceName:   "University House",
					Images: []PostImage{
						{
							ImageURL:  "https://i.imgur.com/NyQRLkf.jpeg",
							Caption:   "View from my room",
							SortOrder: 0,
						},
						{
							ImageURL:  "https://i.imgur.com/4am0zef.jpeg",
							Caption:   "Common area",
							SortOrder: 1,
						},
					},
					Comments: []Comment{
						{
							Username: "alice123",
							Content:  "Really? I hate uhouse",
							Replies: []Comment{
								{
									Username: "raquentin",
									Content:  "k",
								},
							},
						},
					},
					CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
				},
			},
		},
		{
			Email:       "alice@gmail.com",
			Username:    "alice123",
			DisplayName: "Alice",
			Bio:         "asdf",
			Location:    "culc",
			PfpURL:      "https://i.pinimg.com/736x/47/d5/3b/47d53b895b64497082800efd19950c5b.jpg",
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			Posts: []Post{
				{
					Description: "I love gt!",
					PlaceName:   "Bobby-Dodd Stadium",
					Images: []PostImage{
						{
							ImageURL:  "https://i.imgur.com/seOZ2DC.jpeg",
							Caption:   "nachos crushed these ong",
							SortOrder: 0,
						},
					},
					Comments: []Comment{
						{
							Username: "raquentin",
							Content:  "Pleasant",
							Replies: []Comment{
								{
									Username: "alice123",
									Content:  "Ty!",
								},
							},
						},
					},
					CreatedAt: time.Now().Add(-3 * 24 * time.Hour),
				},
			},
		},
		{
			Email:       "buzz@gmail.com",
			Username:    "buzz",
			DisplayName: "Buzz",
			Bio:         "roll bees",
			Location:    "behind you",
			PfpURL:      "https://i.pinimg.com/736x/55/ae/fb/55aefb14f665cc1bb4b28298833f7cce.jpg",
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			Posts: []Post{
				{
					Description: "Great game at the stadium!",
					PlaceName:   "Bobby-Dodd Stadium",
					Images: []PostImage{
						{
							ImageURL:  "https://cdn.prod.website-files.com/5fd923d0e6a54cf4a0c2acce/68b9d52524ebd85dba96122f_Georgia_Tech-MAIN_i.jpg",
							Caption:   "amazing atmosphere",
							SortOrder: 0,
						},
					},
					Comments: []Comment{
						{
							Username: "raquentin",
							Content:  "Go Jackets!",
							Replies: []Comment{
								{
									Username: "buzz",
									Content:  "Buzz buzz!",
								},
							},
						},
					},
					CreatedAt: time.Now().Add(-2 * 24 * time.Hour),
				},
				{
					Description: "Nice study spot at uhouse",
					PlaceName:   "University House",
					Images: []PostImage{
						{
							ImageURL:  "https://i.imgur.com/zGLrTIt.jpeg",
							Caption:   "quiet and peaceful",
							SortOrder: 0,
						},
						{
							ImageURL:  "https://i.imgur.com/dJPSrp3.jpeg",
							Caption:   "quiet and peaceful",
							SortOrder: 0,
						},
					},
					Comments: []Comment{
						{
							Username: "alice123",
							Content:  "I should check this out",
							Replies: []Comment{
								{
									Username: "buzz",
									Content:  "Definitely recommend!",
								},
							},
						},
					},
					CreatedAt: time.Now().Add(-1 * 24 * time.Hour),
				},
			},
		},
	}
}

func (s *Seeder) getPlaceProfiles() []PlaceProfile {
	return []PlaceProfile{
		{
			Name:       "University House",
			Geometry:   "POINT(33.780060 -84.389709)",
			Properties: map[string]any{},
		},
		{
			Name:       "Bobby-Dodd Stadium",
			Geometry:   "POINT(33.7725 -84.392778)",
			Properties: map[string]any{},
		},
	}
}

func (s *Seeder) getFollowRelationships() []FollowRelationship {
	return []FollowRelationship{
		{
			FollowerUsername:  "raquentin",
			FollowingUsername: "alice123",
		},
		{
			FollowerUsername:  "alice123",
			FollowingUsername: "raquentin",
		},
		{
			FollowerUsername:  "raquentin",
			FollowingUsername: "buzz",
		},
		{
			FollowerUsername:  "alice123",
			FollowingUsername: "buzz",
		},
	}
}

func (s *Seeder) seedPlaces(ctx context.Context, placeProfiles []PlaceProfile) ([]*models.Place, error) {
	var places []*models.Place

	for _, profile := range placeProfiles {
		place := &models.Place{
			ID:         uuid.New(),
			Name:       profile.Name,
			Geometry:   stringPtr(profile.Geometry),
			Properties: profile.Properties,
			CreatedAt:  time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:  time.Now().Add(-30 * 24 * time.Hour),
		}

		if err := s.placeRepo.Create(ctx, place); err != nil {
			return nil, fmt.Errorf("failed to create place %s: %w", place.Name, err)
		}

		places = append(places, place)
	}

	return places, nil
}

func (s *Seeder) seedUserProfiles(ctx context.Context, userProfiles []UserProfile, places []*models.Place) ([]*models.User, error) {
	var users []*models.User
	usernameToUser := make(map[string]*models.User)

	for _, profile := range userProfiles {
		existingUser, err := s.userRepo.GetByEmail(ctx, profile.Email)
		if err == nil && existingUser != nil {
			users = append(users, existingUser)
			usernameToUser[profile.Username] = existingUser
			fmt.Printf("UserID for existing user %s: %s\n", profile.Username, existingUser.ID)
			continue
		}

		user := &models.User{
			ID:          uuid.New(),
			Email:       profile.Email,
			Username:    stringPtr(profile.Username),
			DisplayName: stringPtr(profile.DisplayName),
			Bio:         stringPtr(profile.Bio),
			Location:    stringPtr(profile.Location),
			PfpURL:      stringPtr(profile.PfpURL),
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.CreatedAt,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", profile.Email, err)
		}

		users = append(users, user)
		usernameToUser[profile.Username] = user
		fmt.Printf("UserID for new user %s: %s\n", profile.Username, user.ID)
	}

	for _, profile := range userProfiles {
		user := usernameToUser[profile.Username]

		// Create a map of place names to place IDs for easy lookup
		placeNameToID := make(map[string]uuid.UUID)
		for _, place := range places {
			placeNameToID[place.Name] = place.ID
		}

		for _, postData := range profile.Posts {
			placeID, exists := placeNameToID[postData.PlaceName]
			if !exists {
				return nil, fmt.Errorf("place %s not found for post", postData.PlaceName)
			}

			post := &models.Post{
				ID:          uuid.New(),
				UserID:      user.ID,
				PlaceID:     placeID,
				Description: stringPtr(postData.Description),
				CreatedAt:   postData.CreatedAt,
				UpdatedAt:   postData.CreatedAt,
			}

			if err := s.postRepo.Create(ctx, post); err != nil {
				return nil, fmt.Errorf("failed to create post: %w", err)
			}

			for _, imageData := range postData.Images {
				image := &models.PostImage{
					ID:        uuid.New(),
					PostID:    post.ID,
					ImageURL:  imageData.ImageURL,
					Caption:   stringPtr(imageData.Caption),
					SortOrder: imageData.SortOrder,
					CreatedAt: post.CreatedAt,
					UpdatedAt: post.CreatedAt,
				}
				if err := s.postRepo.CreateImage(ctx, image); err != nil {
					return nil, fmt.Errorf("failed to create post image: %w", err)
				}
			}

			for _, commentData := range postData.Comments {
				if err := s.createCommentWithReplies(ctx, post.ID, commentData, usernameToUser); err != nil {
					return nil, fmt.Errorf("failed to create comment: %w", err)
				}
			}
		}
	}

	return users, nil
}

func (s *Seeder) createCommentWithReplies(ctx context.Context, postID uuid.UUID, commentData Comment, usernameToUser map[string]*models.User) error {
	commenter, exists := usernameToUser[commentData.Username]
	if !exists {
		return fmt.Errorf("user %s not found for comment", commentData.Username)
	}

	comment := &models.Comment{
		ID:        uuid.New(),
		PostID:    postID,
		UserID:    commenter.ID,
		Content:   commentData.Content,
		CreatedAt: time.Now().Add(-time.Duration(12) * time.Hour),
		UpdatedAt: time.Now().Add(-time.Duration(12) * time.Hour),
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	for _, replyData := range commentData.Replies {
		replier, exists := usernameToUser[replyData.Username]
		if !exists {
			return fmt.Errorf("user %s not found for reply", replyData.Username)
		}

		reply := &models.Comment{
			ID:        uuid.New(),
			PostID:    postID,
			ParentID:  &comment.ID,
			UserID:    replier.ID,
			Content:   replyData.Content,
			CreatedAt: time.Now().Add(-time.Duration(6) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(6) * time.Hour),
		}

		if err := s.commentRepo.Create(ctx, reply); err != nil {
			return fmt.Errorf("failed to create reply: %w", err)
		}
	}

	return nil
}

func (s *Seeder) seedFollowRelationships(ctx context.Context, users []*models.User, followRelationships []FollowRelationship) error {
	usernameToUser := make(map[string]*models.User)
	for _, user := range users {
		if user.Username != nil {
			usernameToUser[*user.Username] = user
		}
	}

	for _, followData := range followRelationships {
		follower, exists := usernameToUser[followData.FollowerUsername]
		if !exists {
			continue
		}

		following, exists := usernameToUser[followData.FollowingUsername]
		if !exists {
			continue
		}

		follow := &models.Follow{
			ID:          uuid.New(),
			FollowerID:  follower.ID,
			FollowingID: following.ID,
			CreatedAt:   time.Now().Add(-time.Duration(15) * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-time.Duration(15) * 24 * time.Hour),
		}

		if err := s.followRepo.CreateFollow(ctx, follow); err != nil {
			return fmt.Errorf("failed to create follow relationship: %w", err)
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
