package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/pin-app/pin/internal/handlers"
	"github.com/pin-app/pin/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New()

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
