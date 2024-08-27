package apiserver

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/config"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/db"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/mysqldb"
	"github.com/go-sql-driver/mysql"
)

func Run() {
	u1 := &rlnotif.User{
		Name:      "greg",
		Email:     "greg@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cfg := mysql.Config{
		User:                 config.Envs.DBUSer,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := db.NewMySQLStorage(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	// Will be necessary for the API when there is a server
	initStorage(db)

	// User Services test
	userService := mysqldb.NewUserService(db)

	userID, err := userService.CreateUser(u1)
	if err != nil {
		log.Fatal(err)
	}

	u, err := userService.User(userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User found: %v\n\n", u)

	users, err := userService.Users()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range users {
		fmt.Printf("Value: %+v. Type: %T\n\n", v, v)
	}

	// Notification Services test
	n1 := &rlnotif.Notification{
		NotificationType: "Marketing",
		Message:          "Hello",
		UserID:           1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	notificationService := mysqldb.NewNotificationService(db)

	notificationID, err := notificationService.CreateNotification(n1)
	if err != nil {
		log.Fatal(err)
	}

	n, err := notificationService.Notification(notificationID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Notification found: %v\n\n", n)

	notifications, err := notificationService.Notifications()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range notifications {
		fmt.Printf("Value: %+v. Type: %T\n\n", v, v)
	}

	// RateLimit Services test
	rl1 := &rlnotif.RateLimit{
		NotificationType: "Status Update",
		TimeWindow:       "Minute",
		MaxLimit:         2,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rateLimitService := mysqldb.NewRateLimitService(db)

	rateLimitID, err := rateLimitService.CreateRateLimit(rl1)
	if err != nil {
		log.Fatal(err)
	}

	rl, err := rateLimitService.RateLimit(rateLimitID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RateLimit found: %v\n\n", rl)

	rateLimits, err := rateLimitService.RateLimits()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range rateLimits {
		fmt.Printf("Value: %+v. Type: %T\n\n", v, v)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
