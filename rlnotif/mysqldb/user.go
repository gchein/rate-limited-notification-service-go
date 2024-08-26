package mysqldb

import (
	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
)

type UserService struct {
	// db *DB
	DB *UserRepo
}

type UserRepo []*rlnotif.User

// Initial test without actual database
// func NewUserService(db *DB) *UserService {
func NewUserService(db *UserRepo) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) User(id int) (*rlnotif.User, error) {
	user := (*s.DB)[id-1]

	return user, nil
}

func (s *UserService) CreateUser(user *rlnotif.User) error {
	*s.DB = append(*s.DB, user)

	return nil
}
