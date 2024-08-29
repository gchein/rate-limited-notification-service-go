package http

import (
	"net/http"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/utils"
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
	reqBody := rlnotif.Notification{}
	if err := utils.ParseJSON(r, &reqBody); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	errChan := make(chan error)

	go func() {
		errChan <- h.service.Send(reqBody.NotificationType, reqBody.UserID, reqBody.Message)
		close(errChan)
	}()
	if err := <-errChan; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteMessage(w, http.StatusOK, reqBody.Message)
}
