package handlers

import (
	"net/http"

	"github.com/pin-app/pin/internal/server"
)

func RegisterRoutes(srv *server.Server) {
	router := srv.GetRouter()

	router.HandleFunc("/health", "GET", HealthCheck(srv))
}
