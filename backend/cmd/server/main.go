package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/pin-app/pin/internal/database"
	"github.com/pin-app/pin/internal/handlers"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/seed"
	"github.com/pin-app/pin/internal/server"
	"github.com/pin-app/pin/migrations"
)

func main() {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		slog.Error("failed to prepare upload directory", "error", err)
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if p, err := strconv.Atoi(port); err != nil || p <= 0 || p > 65535 {
		slog.Error("invalid PORT environment variable; must be 1-65535",
			"PORT", port,
			"error", err,
		)
		os.Exit(1)
	}

	devMode := os.Getenv("DEV_MODE") == "true"
	if devMode {
		slog.Info("running in development mode - authentication bypassed for dev users")
	}

	var db *database.DB
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		slog.Info("running database migrations")
		if err := migrations.Run(dbURL); err != nil {
			slog.Error("database migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("database migrations complete")

		var err error
		db, err = database.New(dbURL)
		if err != nil {
			slog.Error("failed to open database connection", "error", err)
			os.Exit(1)
		}
		defer db.Close()
	}

	if devMode && db == nil {
		slog.Error("dev mode requires DATABASE_URL to be set for seeding dummy data")
		os.Exit(1)
	}

	var srv *server.Server
	if db != nil {
		srv = server.NewWithDB(db.GetConnection())

		if devMode {
			slog.Info("seeding development data")

			userRepo := repository.NewUserRepository(db)
			placeRepo := repository.NewPlaceRepository(db)
			postRepo := repository.NewPostRepository(db)
			commentRepo := repository.NewCommentRepository(db)
			ratingRepo := repository.NewRatingRepository(db)
			followRepo := repository.NewFollowRepository(db)
			likeRepo := repository.NewLikeRepository(db)

			seeder := seed.NewSeeder(
				userRepo,
				placeRepo,
				postRepo,
				ratingRepo,
				commentRepo,
				followRepo,
				likeRepo,
			)

			if err := seeder.SeedDevData(context.Background()); err != nil {
				slog.Error("failed to seed development data", "error", err)
			} else {
				slog.Info("development data seeded successfully")
			}
		}
	} else {
		srv = server.New()
	}

	srv.ServeStatic("/uploads/", uploadDir)

	handlers.RegisterRoutes(srv, db, uploadDir)

	slog.Info("server starting",
		"port", port,
		"service", "pin",
	)

	if err := http.ListenAndServe(":"+port, srv); err != nil {
		slog.Error("server failed to start",
			"error", err,
			"port", port,
		)
		os.Exit(1)
	}
}
