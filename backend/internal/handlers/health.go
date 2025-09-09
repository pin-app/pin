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
	Count     int       `json:"count"`
}

func HealthCheck(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := r.Context().Value("logger").(*slog.Logger)
		if !ok {
			logger = slog.Default()
		}

		if err := srv.StoreHealthCheck(r.RemoteAddr, r.UserAgent()); err != nil {
			logger.Error("failed to store health check", "error", err)
		}

		count, err := srv.GetHealthCheckCount()
		if err != nil {
			logger.Error("failed to get health check count", "error", err)
			count = 0
		}

		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "pin",
			Count:     count,
		}

		logger.Info("health check requested",
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status", response.Status,
			"count", count,
		)

		server.WriteJSON(w, http.StatusOK, response)
	}
}
