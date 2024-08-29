package http

import (
	"log"
	"net/http"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif/db"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/mysqldb"
	"github.com/gorilla/mux"
)

type Server struct {
	addr string
	db   *db.DB
}

func NewServer(addr string, db *db.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()

	notificationService := mysqldb.NewNotificationService(*s.db)
	notificationHandler := NewNotificationHandler(&notificationService)
	notificationHandler.RegisterNotificationRoutes(router)

	rateLimitService := mysqldb.NewRateLimitService(*s.db)
	rateLimitHandler := NewRateLimitHandler(&rateLimitService)
	rateLimitHandler.RegisterRateLimitRoutes(router)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
