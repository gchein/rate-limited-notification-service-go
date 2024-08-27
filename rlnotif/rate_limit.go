package rlnotif

import (
	"time"
)

type RateLimit struct {
	ID               int64  `json:"id"`
	NotificationType string `json:"notificationType"`
	TimeWindow       string `json:"timeWindow"`
	MaxLimit         int    `json:"maxLimit"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RateLimitService interface {
	RateLimit(id int64) (*RateLimit, error)
	RateLimits() ([]*RateLimit, error)
	CreateRateLimit(notif *RateLimit) error
}
