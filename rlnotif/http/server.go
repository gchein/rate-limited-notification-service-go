package http

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif/mysqldb"
	"github.com/gorilla/mux"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()

	userService := mysqldb.NewUserService(s.db)
	userHandler := NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	// [...]

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
