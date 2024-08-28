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
	RateLimitStorage
	RateLimitCacher
}

type RateLimitStorage interface {
	RateLimits() ([]*RateLimit, error)
	CreateRateLimit(rateLimit *RateLimit) (int64, error)
	DeleteRateLimit(id int64) error
}

type RateLimitCacher interface {
	UpdateRateLimitsCache() error
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
