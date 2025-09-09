package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pin-app/pin/internal/handlers"
	"github.com/pin-app/pin/internal/server"
	"github.com/pin-app/pin/migrations"
)

func main() {
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

	var db *sql.DB
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		slog.Info("running database migrations")
		if err := migrations.Run(dbURL); err != nil {
			slog.Error("database migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("database migrations complete")

		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			slog.Error("failed to open database connection", "error", err)
			os.Exit(1)
		}
		defer db.Close()
	}

	var srv *server.Server
	if db != nil {
		srv = server.NewWithDB(db)
	} else {
		srv = server.New()
	}

	handlers.RegisterRoutes(srv)

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
