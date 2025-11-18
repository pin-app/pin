package handlers

import (
	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/middleware"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/server"
)

func RegisterRoutes(srv *server.Server, db *database.DB, uploadDir string) {
	router := srv.GetRouter()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	placeRepo := repository.NewPlaceRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	ratingRepo := repository.NewRatingRepository(db)
	oauthRepo := repository.NewOAuthRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	followRepo := repository.NewFollowRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// Initialize auth middleware
	authMW := middleware.NewAuthMiddleware(sessionRepo, userRepo)

	// Initialize handlers
	userHandler := NewUserHandler(userRepo)
	placeHandler := NewPlaceHandler(placeRepo)
	postHandler := NewPostHandler(postRepo, placeRepo, userRepo, commentRepo, likeRepo)
	commentHandler := NewCommentHandler(commentRepo, postRepo, userRepo)
	uploadHandler := NewUploadHandler(uploadDir)
	ratingHandler := NewRatingHandler(ratingRepo, placeRepo, userRepo)
	oauthHandler := NewOAuthHandler(oauthRepo, userRepo, sessionRepo)
	followHandler := NewFollowHandler(followRepo, userRepo)

	// OAuth routes (public)
	// Upload routes
	router.HandleFunc("/api/uploads", "POST", authMW.RequireAuth(uploadHandler.UploadImage))

	router.HandleFunc("/api/auth/google", "GET", oauthHandler.GoogleAuth)
	router.HandleFunc("/api/auth/google/callback", "GET", oauthHandler.GoogleCallback)
	router.HandleFunc("/api/auth/apple", "GET", oauthHandler.AppleAuth)
	router.HandleFunc("/api/auth/apple/callback", "GET", oauthHandler.AppleCallback)
	router.HandleFunc("/api/auth/logout", "POST", oauthHandler.Logout)

	// User routes
	router.HandleFunc("/api/users", "POST", userHandler.CreateUser)
	router.HandleFunc("/api/users", "GET", authMW.OptionalAuth(userHandler.ListUsers))
	router.HandleFunc("/api/users/search", "GET", authMW.OptionalAuth(userHandler.SearchUsers))
	router.HandleFunc("/api/users/{id}", "GET", authMW.OptionalAuth(userHandler.GetUser))
	router.HandleFunc("/api/users/{id}", "PUT", authMW.RequireAuth(userHandler.UpdateUser))
	router.HandleFunc("/api/users/{id}", "DELETE", authMW.RequireAuth(userHandler.DeleteUser))

	// Follow routes
	router.HandleFunc("/api/users/{id}/follow", "POST", authMW.RequireAuth(followHandler.FollowUser))
	router.HandleFunc("/api/users/{id}/follow", "DELETE", authMW.RequireAuth(followHandler.UnfollowUser))
	router.HandleFunc("/api/users/{id}/following", "GET", authMW.OptionalAuth(followHandler.GetFollowing))
	router.HandleFunc("/api/users/{id}/followers", "GET", authMW.OptionalAuth(followHandler.GetFollowers))
	router.HandleFunc("/api/users/{id}/follow-status", "GET", authMW.RequireAuth(followHandler.CheckFollowStatus))
	router.HandleFunc("/api/users/{id}/stats", "GET", authMW.OptionalAuth(followHandler.GetUserStats))

	// Place routes
	router.HandleFunc("/api/places", "POST", authMW.RequireAuth(placeHandler.CreatePlace))
	router.HandleFunc("/api/places", "GET", authMW.OptionalAuth(placeHandler.ListPlaces))
	router.HandleFunc("/api/places/search", "GET", authMW.OptionalAuth(placeHandler.SearchPlaces))
	router.HandleFunc("/api/places/nearby", "GET", authMW.OptionalAuth(placeHandler.SearchNearbyPlaces))
	router.HandleFunc("/api/places/{id}", "GET", authMW.OptionalAuth(placeHandler.GetPlace))
	router.HandleFunc("/api/places/{id}", "PUT", authMW.RequireAuth(placeHandler.UpdatePlace))
	router.HandleFunc("/api/places/{id}", "DELETE", authMW.RequireAuth(placeHandler.DeletePlace))

	// Post routes
	router.HandleFunc("/api/posts", "POST", authMW.RequireAuth(postHandler.CreatePost))
	router.HandleFunc("/api/posts", "GET", authMW.OptionalAuth(postHandler.ListPosts))
	router.HandleFunc("/api/posts/{id}", "GET", authMW.OptionalAuth(postHandler.GetPost))
	router.HandleFunc("/api/posts/{id}", "PUT", authMW.RequireAuth(postHandler.UpdatePost))
	router.HandleFunc("/api/posts/{id}", "DELETE", authMW.RequireAuth(postHandler.DeletePost))
	router.HandleFunc("/api/posts/{id}/likes", "POST", authMW.RequireAuth(postHandler.LikePost))
	router.HandleFunc("/api/posts/{id}/likes", "DELETE", authMW.RequireAuth(postHandler.UnlikePost))
	router.HandleFunc("/api/users/{id}/posts", "GET", authMW.OptionalAuth(postHandler.ListPostsByUser))
	router.HandleFunc("/api/places/{id}/posts", "GET", authMW.OptionalAuth(postHandler.ListPostsByPlace))

	// Comment routes
	router.HandleFunc("/api/comments", "POST", authMW.RequireAuth(commentHandler.CreateComment))
	router.HandleFunc("/api/comments/{id}", "GET", authMW.OptionalAuth(commentHandler.GetComment))
	router.HandleFunc("/api/comments/{id}", "PUT", authMW.RequireAuth(commentHandler.UpdateComment))
	router.HandleFunc("/api/comments/{id}", "DELETE", authMW.RequireAuth(commentHandler.DeleteComment))
	router.HandleFunc("/api/comments/{id}/replies", "GET", authMW.OptionalAuth(commentHandler.GetCommentReplies))
	router.HandleFunc("/api/posts/{id}/comments", "GET", authMW.OptionalAuth(commentHandler.ListCommentsByPost))
	router.HandleFunc("/api/users/{id}/comments", "GET", authMW.OptionalAuth(commentHandler.ListCommentsByUser))

	// Rating routes
	router.HandleFunc("/api/places/{id}/ratings", "POST", authMW.RequireAuth(ratingHandler.CreateRating))
	router.HandleFunc("/api/places/{id}/ratings", "GET", authMW.OptionalAuth(ratingHandler.ListRatingsByPlace))
	router.HandleFunc("/api/places/{id}/ratings", "PUT", authMW.RequireAuth(ratingHandler.UpdateRating))
	router.HandleFunc("/api/places/{id}/ratings", "DELETE", authMW.RequireAuth(ratingHandler.DeleteRating))
	router.HandleFunc("/api/places/{id}/ratings/me", "GET", authMW.RequireAuth(ratingHandler.GetRating))
	router.HandleFunc("/api/places/{id}/ratings/average", "GET", authMW.OptionalAuth(ratingHandler.GetAverageRating))
	router.HandleFunc("/api/places/compare", "POST", authMW.RequireAuth(ratingHandler.CreateComparison))
	router.HandleFunc("/api/users/{id}/comparisons", "GET", authMW.RequireAuth(ratingHandler.ListComparisonsByUser))
}
