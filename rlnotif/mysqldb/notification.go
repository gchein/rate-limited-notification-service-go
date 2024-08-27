package mysqldb

import (
	"database/sql"
	"fmt"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
)

type NotificationService struct {
	DB *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{DB: db}
}

func (s *NotificationService) Notification(id int64) (*rlnotif.Notification, error) {
	db := s.DB

	var notification rlnotif.Notification

	row := db.QueryRow("SELECT * FROM notifications WHERE id = ?", id)
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
	db := s.DB
	var notifications []*rlnotif.Notification

	rows, err := db.Query("SELECT * FROM notifications")
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

func (s *NotificationService) CreateNotification(notification *rlnotif.Notification) (int64, error) {
	db := s.DB

	result, err := db.Exec("INSERT INTO notifications (notification_type, message, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		&notification.NotificationType,
		&notification.Message,
		&notification.UserID,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateNotification: %v", err)
	}
	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateNotification: %v", err)
	}

	return ID, nil
}
