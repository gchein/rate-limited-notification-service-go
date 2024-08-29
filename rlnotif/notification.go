package rlnotif

import (
	"time"
)

type Notification struct {
	ID               int64  `json:"id"`
	NotificationType string `json:"notificationType"`
	Message          string `json:"message"`
	UserID           int64  `json:"userId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NotificationService interface {
	NotificationStorage
	NotificationSender
}

type NotificationStorage interface {
	Notification(id int64) (*Notification, error)
	Notifications() ([]*Notification, error)
	CreateNotification(notification *Notification) error
}

type NotificationSender interface {
	Send(notificationType string, userId int64, message string) error
}
