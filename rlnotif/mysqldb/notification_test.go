package mysqldb

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/config"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/db"
	"github.com/go-sql-driver/mysql"
)

func initTestStorage() (DB *sql.DB) {
	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.TestDBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}
	DB, err := db.NewMySQLStorage(&cfg)

	if err != nil {
		log.Fatal("Error initializing DB:", err)
		return nil
	}

	return DB
}

func TestCreateNotification(t *testing.T) {
	DB := initTestStorage()
	if DB == nil {
		log.Fatal("Failed to initialize database connection")
	}

	tx, err := DB.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})

	userService := NewUserService(tx)

	u := &rlnotif.User{
		Name:      "Test User",
		Email:     "test_user@mail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uID, err := userService.CreateUser(u)
	if err != nil {
		log.Fatal(err)
	}

	rateLimitService := NewRateLimitService(tx)

	rl := &rlnotif.RateLimit{
		NotificationType: "Test Notification",
		TimeWindow:       "Minute",
		MaxLimit:         1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = rateLimitService.CreateRateLimit(rl)
	if err != nil {
		log.Fatal(err)
	}

	notificationService := NewNotificationService(tx)

	n := &rlnotif.Notification{
		NotificationType: "Test Notification",
		Message:          "Test message",
		UserID:           uID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	got := notificationService.CreateNotification(n)
	if got != nil {
		t.Errorf("Could not create notification")
	}
}

func TestSendNotification(t *testing.T) {
	DB := initTestStorage()
	if DB == nil {
		log.Fatal("Failed to initialize database connection")
	}

	tx, err := DB.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})

	userService := NewUserService(tx)

	u := &rlnotif.User{
		Name:      "Test User",
		Email:     "test_user@mail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uID, err := userService.CreateUser(u)
	if err != nil {
		log.Fatal(err)
	}

	rateLimitService := NewRateLimitService(tx)

	rl := &rlnotif.RateLimit{
		NotificationType: "Test Notification",
		TimeWindow:       "Minute",
		MaxLimit:         1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = rateLimitService.CreateRateLimit(rl)
	if err != nil {
		log.Fatal(err)
	}

	notificationService := NewNotificationService(tx)

	n := &rlnotif.Notification{
		NotificationType: "Test Notification",
		Message:          "Test message",
		UserID:           uID,
	}

	err = notificationService.Send(
		n.NotificationType, n.UserID, n.Message,
	)
	if err != nil {
		log.Fatal(err)
	}

	got := notificationService.Send(
		n.NotificationType, n.UserID, n.Message,
	)
	if got == nil {
		t.Errorf("Send Method disrespected Rate Limit")
	}
}
