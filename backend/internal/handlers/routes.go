package handlers

import (
	"github.com/pin-app/pin/internal/server"
)

func RegisterRoutes(srv *server.Server) {
	router := srv.GetRouter()

	router.HandleFunc("/health", "GET", HealthCheck(srv))
}
