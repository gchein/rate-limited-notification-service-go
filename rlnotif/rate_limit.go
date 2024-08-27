package rlnotif

import (
	"sync"
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
	// RateLimit(id int64) (*RateLimit, error)
	RateLimits() ([]*RateLimit, error)
	// RateLimitsByFields() ([]*RateLimit, error) // Alter to lookup the DB on query fields // Cuidado aqui pra acesso simultaneo à DB que pode estar sendo alterada
	CreateRateLimit(rateLimit *RateLimit) (int64, error)
	// DeleteRateLimit()
}

type rateLimitsCache struct {
	mutex  sync.Mutex
	limits map[string]map[string]int
}

var cachedLimits = &rateLimitsCache{
	limits: make(map[string]map[string]int),
}

func CacheRateLimits(rateLimits []*RateLimit) {
	cachedLimits.mutex.Lock()
	defer cachedLimits.mutex.Unlock()

	newLimits := make(map[string]map[string]int)

	for _, rl := range rateLimits {
		nt := rl.NotificationType
		tw := rl.TimeWindow
		lim := rl.MaxLimit

		if notifMap, exists := newLimits[nt]; exists {
			notifMap[tw] = lim
		} else {
			m := make(map[string]int)
			m[tw] = lim
			newLimits[nt] = m
		}
	}

	cachedLimits.limits = newLimits
}

func RateLimitsFromCache(notificationType string) (map[string]int, bool) {
	cachedLimits.mutex.Lock()
	defer cachedLimits.mutex.Unlock()

	rateLimits, exists := cachedLimits.limits[notificationType]

	return rateLimits, exists
}
