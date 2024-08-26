package rlnotif

import (
	"time"
)

type Notification struct {
	ID               int    `json:"id"`
	NotificationType string `json:"notificationType"`
	Message          string `json:"message"`
	UserID           int    `json:"userId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NotificationService interface {
	Notification(id int) (*Notification, error)
	CreateNotification(notif *Notification) error
}
