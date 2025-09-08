package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/pin-app/pin/internal/server"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Get logger from context
	logger, ok := r.Context().Value("logger").(*slog.Logger)
	if !ok {
		// Fallback to default logger if not found
		logger = slog.Default()
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "pin",
	}

	logger.Info("health check requested",
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
		"status", response.Status,
	)

	server.WriteJSON(w, http.StatusOK, response)
}
