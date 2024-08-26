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
	u3 := &rlnotif.User{
		ID:        3,
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

	s := mysqldb.NewUserService(db)

	userCreateErr := s.CreateUser(u3)
	if userCreateErr != nil {
		log.Fatal(userCreateErr)
	}

	u, err := s.User(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User found: %v\n\n", u)

	users, err := s.Users()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range users {
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
