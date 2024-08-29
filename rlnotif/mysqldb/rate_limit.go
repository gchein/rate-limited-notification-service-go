package mysqldb

import (
	"fmt"
	"regexp"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/db"
)

type RateLimitService struct {
	DB db.DB
}

func NewRateLimitService(db db.DB) RateLimitService {
	return RateLimitService{DB: db}
}

// Ensure service implements interface.
var _ rlnotif.RateLimitService = (*RateLimitService)(nil)

func (s *RateLimitService) RateLimits() ([]*rlnotif.RateLimit, error) {
	var rateLimits []*rlnotif.RateLimit

	rows, err := s.DB.Query("SELECT * FROM rate_limits ORDER BY notification_type, max_limit")
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
	result, err := s.DB.Exec("INSERT INTO rate_limits (notification_type, time_window, max_limit, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		&rateLimit.NotificationType,
		&rateLimit.TimeWindow,
		&rateLimit.MaxLimit,
		&rateLimit.CreatedAt,
		&rateLimit.UpdatedAt,
	)
	if err != nil {
		re := regexp.MustCompile(`Duplicate entry`)
		if re.MatchString(err.Error()) {
			return 0, fmt.Errorf("CreateRateLimit: Rate Limit already implemented for notification type '%v' and time window '%v'",
				rateLimit.NotificationType,
				rateLimit.TimeWindow)
		}

		return 0, fmt.Errorf("CreateRateLimit: %v", err)
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateRateLimit: %v", err)
	}

	err = s.UpdateRateLimitsCache()
	if err != nil {
		return ID, err
	}

	return ID, nil
}

func (s *RateLimitService) DeleteRateLimit(id int64) error {
	query := "DELETE FROM rate_limits WHERE id = ?"

	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete rate limit: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not determine the number of rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rate limit found with id %d", id)
	}

	err = s.UpdateRateLimitsCache()
	if err != nil {
		return err
	}

	return nil
}

func (s *RateLimitService) UpdateRateLimitsCache() error {
	newRateLimits, err := s.RateLimits()
	if err != nil {
		return fmt.Errorf("RateLimits: %v", err)
	}

	rlnotif.CacheRateLimits(newRateLimits)

	return nil
}
