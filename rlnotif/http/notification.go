package http

import (
	"net/http"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	service rlnotif.NotificationService // Interface Type
}

func NewNotificationHandler(service rlnotif.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) RegisterNotificationRoutes(router *mux.Router) {
	router.HandleFunc("/notifications", h.handlePostNotifications).Methods(http.MethodPost)
}

func (h *NotificationHandler) handlePostNotifications(w http.ResponseWriter, r *http.Request) {
	// h.service.Send()

	// notifications, err := h.service.Notifications()
	// if err != nil {
	// 	// utils.WriteError(w, http.StatusInternalServerError, err)
	// }

	// utils.WriteJSON(w, http.StatusOK, notifications)
}
