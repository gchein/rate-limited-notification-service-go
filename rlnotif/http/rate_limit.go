package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/utils"
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
	router.HandleFunc("/rate-limits/{ID}", h.handleDeleteRateLimit).Methods(http.MethodDelete)
}

func (h *RateLimitHandler) handleGetRateLimits(w http.ResponseWriter, r *http.Request) {
	rateLimits, err := h.service.RateLimits()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, rateLimits)
}

func (h *RateLimitHandler) handlePostRateLimits(w http.ResponseWriter, r *http.Request) {
	rl := rlnotif.RateLimit{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := utils.ParseJSON(r, &rl); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.CreateRateLimit(&rl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	rl.ID = id
	utils.WriteJSON(w, http.StatusCreated, rl)
}

func (h *RateLimitHandler) handleDeleteRateLimit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rateLimitId, ok := vars["ID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID not found on request"))
		return
	}

	rateLimitID, err := strconv.Atoi(rateLimitId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid ID"))
		return
	}

	err = h.service.DeleteRateLimit(int64(rateLimitID))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	utils.WriteMessage(w, http.StatusOK, "rate limit successfully deleted")
}
