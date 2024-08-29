package apiserver

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/config"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/db"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/http"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/mysqldb"
	"github.com/go-sql-driver/mysql"
)

func Run() {
	db := initStorage()

	rateLimitService := mysqldb.NewRateLimitService(db)

	rateLimits, err := rateLimitService.RateLimits()
	if err != nil {
		log.Fatal(err)
	}
	go rlnotif.CacheRateLimits(rateLimits)

	server := http.NewServer(fmt.Sprintf(":%s", config.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage() (DB *sql.DB) {
	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}
	DB, err := db.NewMySQLStorage(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")

	return DB
}

func Seed() {
	db := initStorage()

	// User seeds
	log.Println("Seeding users...")
	userService := mysqldb.NewUserService(db)

	u1 := &rlnotif.User{
		Name:      "greg",
		Email:     "greg@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	u2 := &rlnotif.User{
		Name:      "suzy",
		Email:     "suzy@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	u3 := &rlnotif.User{
		Name:      "john",
		Email:     "john@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users := []*rlnotif.User{
		u1,
		u2,
		u3,
	}

	for _, v := range users {
		_, err := userService.CreateUser(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Rate Limit seeds
	log.Println("Seeding rate limits...")
	rateLimitService := mysqldb.NewRateLimitService(db)

	rl1 := &rlnotif.RateLimit{
		NotificationType: "Status Update",
		TimeWindow:       "Minute",
		MaxLimit:         2,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl2 := &rlnotif.RateLimit{
		NotificationType: "Status Update",
		TimeWindow:       "Hour",
		MaxLimit:         5,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl3 := &rlnotif.RateLimit{
		NotificationType: "Status Update",
		TimeWindow:       "Day",
		MaxLimit:         20,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl4 := &rlnotif.RateLimit{
		NotificationType: "Daily News",
		TimeWindow:       "Day",
		MaxLimit:         1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl5 := &rlnotif.RateLimit{
		NotificationType: "Marketing",
		TimeWindow:       "Hour",
		MaxLimit:         3,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl6 := &rlnotif.RateLimit{
		NotificationType: "Marketing",
		TimeWindow:       "Day",
		MaxLimit:         10,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl7 := &rlnotif.RateLimit{
		NotificationType: "Project Invitation",
		TimeWindow:       "Day",
		MaxLimit:         1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rl8 := &rlnotif.RateLimit{
		NotificationType: "Project Invitation",
		TimeWindow:       "Month",
		MaxLimit:         3,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	rate_limits := []*rlnotif.RateLimit{
		rl1,
		rl2,
		rl3,
		rl4,
		rl5,
		rl6,
		rl7,
		rl8,
	}

	for _, v := range rate_limits {
		_, err := rateLimitService.CreateRateLimit(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}
