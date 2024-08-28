package mysqldb

import (
	"database/sql"
	"fmt"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
)

type RateLimitService struct {
	DB *sql.DB
}

func NewRateLimitService(db *sql.DB) *RateLimitService {
	return &RateLimitService{DB: db}
}

// Ensure service implements interface.
var _ rlnotif.RateLimitService = (*RateLimitService)(nil)

func (s *RateLimitService) RateLimit(id int64) (*rlnotif.RateLimit, error) {
	db := s.DB

	var rateLimit rlnotif.RateLimit

	row := db.QueryRow("SELECT * FROM rate_limits WHERE id = ?", id)
	if err := row.Scan(
		&rateLimit.ID,
		&rateLimit.NotificationType,
		&rateLimit.TimeWindow,
		&rateLimit.MaxLimit,
		&rateLimit.CreatedAt,
		&rateLimit.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &rateLimit, fmt.Errorf("RateLimit %d: no such rate limit", id)
		}
		return &rateLimit, fmt.Errorf("RateLimit %d: %v", id, err)
	}
	return &rateLimit, nil
}

func (s *RateLimitService) RateLimits() ([]*rlnotif.RateLimit, error) {
	db := s.DB
	var rateLimits []*rlnotif.RateLimit

	rows, err := db.Query("SELECT * FROM rate_limits ORDER BY notification_type, max_limit")
	if err != nil {
		return nil, fmt.Errorf("RateLimits: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rateLimit rlnotif.RateLimit
		if err := rows.Scan(
			&rateLimit.ID,
			&rateLimit.NotificationType,
			&rateLimit.TimeWindow,
			&rateLimit.MaxLimit,
			&rateLimit.CreatedAt,
			&rateLimit.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("RateLimits: %v", err)
		}
		rateLimits = append(rateLimits, &rateLimit)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RateLimits: %v", err)
	}
	return rateLimits, nil
}

func (s *RateLimitService) CreateRateLimit(rateLimit *rlnotif.RateLimit) (int64, error) {
	db := s.DB

	result, err := db.Exec("INSERT INTO rate_limits (notification_type, time_window, max_limit, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		&rateLimit.NotificationType,
		&rateLimit.TimeWindow,
		&rateLimit.MaxLimit,
		&rateLimit.CreatedAt,
		&rateLimit.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateRateLimit: %v", err)
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateRateLimit: %v", err)
	}

	err = UpdateRateLimitsCache(s)
	if err != nil {
		return ID, err
	}

	return ID, nil
}

func UpdateRateLimitsCache(s *RateLimitService) error {
	newRateLimits, err := s.RateLimits()
	if err != nil {
		return fmt.Errorf("RateLimits: %v", err)
	}

	rlnotif.CacheRateLimits(newRateLimits)

	return nil
}
