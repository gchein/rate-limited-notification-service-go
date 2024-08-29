package http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/jsonutil"
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
	if err := jsonutil.ParseJSON(r, &reqBody); err != nil {
		re := regexp.MustCompile(`EOF`)
		if re.MatchString(err.Error()) {
			jsonutil.WriteError(w, http.StatusBadRequest, fmt.Errorf("please send a valid request body"))
			return
		}

		jsonutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := checkNotificationParams(&reqBody)
	if err != nil {
		jsonutil.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.service.Send(reqBody.NotificationType, reqBody.UserID, reqBody.Message); err != nil {
		jsonutil.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	jsonutil.WriteMessage(w, http.StatusOK, reqBody.Message)
}

func checkNotificationParams(rl *rlnotif.Notification) error {
	if rl.NotificationType == "" {
		return fmt.Errorf("please provide a valid notificationType")
	}

	if rl.UserID <= 0 {
		return fmt.Errorf("please provide a valid userId")
	}

	if rl.Message == "" {
		return fmt.Errorf("please provide a message that is not empty")
	}

	return nil
}
