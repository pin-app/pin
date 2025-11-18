package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/pin-app/pin/internal/middleware"
	"github.com/pin-app/pin/internal/models"
	"github.com/pin-app/pin/internal/repository"
	"github.com/pin-app/pin/internal/server"
)

type NotificationHandler struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationHandler(notificationRepo repository.NotificationRepository, userRepo repository.UserRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	notifications, err := h.notificationRepo.ListByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list notifications"})
		return
	}

	actorCache := make(map[uuid.UUID]*models.UserResponse)
	responses := make([]models.NotificationResponse, len(notifications))

	for i, notification := range notifications {
		resp := notification.ToResponse()
		if actorResp, ok := actorCache[notification.ActorID]; ok {
			resp.Actor = actorResp
		} else {
			actor, err := h.userRepo.GetByID(r.Context(), notification.ActorID)
			if err == nil {
				actorResp := actor.ToResponse()
				resp.Actor = &actorResp
				actorCache[notification.ActorID] = &actorResp
			}
		}
		responses[i] = resp
	}

	server.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"notifications": responses,
		"limit":         limit,
		"offset":        offset,
		"count":         len(responses),
	})
}

func (h *NotificationHandler) ClearNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
		return
	}

	if err := h.notificationRepo.ClearByUserID(r.Context(), userID); err != nil {
		server.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to clear notifications"})
		return
	}

	server.WriteJSON(w, http.StatusOK, map[string]string{"status": "cleared"})
}
