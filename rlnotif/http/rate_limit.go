package http

import (
	"net/http"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gorilla/mux"
)

type RateLimitHandler struct {
	service rlnotif.RateLimitService // Interface Type
}

func NewRateLimitHandler(service rlnotif.RateLimitService) *RateLimitHandler {
	return &RateLimitHandler{service: service}
}

func (h *RateLimitHandler) RegisterRateLimitRoutes(router *mux.Router) {
	router.HandleFunc("/rate-limits", h.handleGetRateLimits).Methods(http.MethodGet)
	router.HandleFunc("/rate-limits", h.handlePostRateLimits).Methods(http.MethodPost)
	router.HandleFunc("/rate-limits", h.handleDeleteRateLimit).Methods(http.MethodDelete)
}

func (h *RateLimitHandler) handleGetRateLimits(w http.ResponseWriter, r *http.Request) {
	// // h.service.RateLimits()

	// rateLimits, err := h.service.RateLimits()
	// if err != nil {
	// 	// utils.WriteError(w, http.StatusInternalServerError, err)
	// }

	// utils.WriteJSON(w, http.StatusOK, rateLimits)
}

func (h *RateLimitHandler) handlePostRateLimits(w http.ResponseWriter, r *http.Request) {
	// h.service.CreateRateLimits()

	// rateLimits, err := h.service.RateLimits()
	// if err != nil {
	// 	// utils.WriteError(w, http.StatusInternalServerError, err)
	// }

	// utils.WriteJSON(w, http.StatusOK, rateLimits)
}

func (h *RateLimitHandler) handleDeleteRateLimit(w http.ResponseWriter, r *http.Request) {
	// h.service.DeleteRateLimits()

	// rateLimits, err := h.service.RateLimits()
	// if err != nil {
	// 	// utils.WriteError(w, http.StatusInternalServerError, err)
	// }

	// utils.WriteJSON(w, http.StatusOK, rateLimits)
}
