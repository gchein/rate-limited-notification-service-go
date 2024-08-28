package http

import (
	"fmt"
	"net/http"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gorilla/mux"
)

type Handler struct {
	service rlnotif.UserService // Interface Type
}

func NewUserHandler(service rlnotif.UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.handleGetUsers).Methods(http.MethodGet)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.Users()
	if err != nil {
		// utils.WriteError(w, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	for _, u := range users {
		fmt.Println(*u)
	}

	// utils.WriteJSON(w, http.StatusOK, users)
}
