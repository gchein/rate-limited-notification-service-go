package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
)

type NotificationService struct {
	DB *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{DB: db}
}

// Ensure service implements interface.
var _ rlnotif.NotificationService = (*NotificationService)(nil)

func (s *NotificationService) Notification(id int64) (*rlnotif.Notification, error) {
	var notification rlnotif.Notification

	row := s.DB.QueryRow("SELECT * FROM notifications WHERE id = ?", id)
	if err := row.Scan(
		&notification.ID,
		&notification.NotificationType,
		&notification.Message,
		&notification.UserID,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &notification, fmt.Errorf("Notification %d: no such notification", id)
		}
		return &notification, fmt.Errorf("Notification %d: %v", id, err)
	}
	return &notification, nil
}

func (s *NotificationService) Notifications() ([]*rlnotif.Notification, error) {
	var notifications []*rlnotif.Notification

	rows, err := s.DB.Query("SELECT * FROM notifications")
	if err != nil {
		return nil, fmt.Errorf("Notifications: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification rlnotif.Notification
		if err := rows.Scan(
			&notification.ID,
			&notification.NotificationType,
			&notification.Message,
			&notification.UserID,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Notifications: %v", err)
		}
		notifications = append(notifications, &notification)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Notifications: %v", err)
	}
	return notifications, nil
}

func (s *NotificationService) CreateNotification(notification *rlnotif.Notification) error {
	_, err := s.DB.Exec("INSERT INTO notifications (notification_type, message, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		&notification.NotificationType,
		&notification.Message,
		&notification.UserID,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("createNotification: %v", err)
	}

	return nil
}

func (s *NotificationService) Send(notificationType string, userID int64, message string) error {
	err := canSendToUser(s, notificationType, userID)
	if err != nil {
		return err
	}

	n := &rlnotif.Notification{
		NotificationType: notificationType,
		Message:          message,
		UserID:           userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.CreateNotification(n)
	if err != nil {
		return fmt.Errorf("Send: %v", err)
	}

	return nil
}

func canSendToUser(s *NotificationService, notificationType string, userID int64) error {

	rateLimitsPerType, exists := rlnotif.RateLimitsFromCache(notificationType)
	if !exists {
		return fmt.Errorf("please verify the Notification Type provided")
	}

	query := `
		SELECT
					notification_type
	`
	var limits []int

	for tw, lim := range rateLimitsPerType {
		limits = append(limits, lim)

		query += fmt.Sprintf(
			", SUM(CASE WHEN created_at >= NOW() - INTERVAL 1 %s THEN 1 ELSE 0 END) AS count_last_%s",
			strings.ToUpper(tw), strings.ToLower(tw))
	}

	query += `
		FROM
				notifications
		WHERE
				user_id = ?
				AND notification_type = ?
		GROUP BY
				notification_type;
	`
	numOfTimeWindows := len(limits)
	notifCountByTimeWindow := make([]int, numOfTimeWindows)
	var scanNotifType string
	scanResult := make([]interface{}, numOfTimeWindows+1)

	scanResult[0] = &scanNotifType
	for i := 1; i <= numOfTimeWindows; i++ {
		scanResult[i] = &notifCountByTimeWindow[i-1]
	}

	row := s.DB.QueryRow(query, userID, notificationType)
	if err := row.Scan(scanResult...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("error fetching notifications on the database for user_id %v, notification type '%v': %v",
			userID,
			notificationType,
			err,
		)
	}

	for i, count := range notifCountByTimeWindow {
		if count == limits[i] {
			return fmt.Errorf("max Notification Limit reached for user_id %v, notification type '%v'",
				userID,
				notificationType)
		}
	}

	return nil
}
