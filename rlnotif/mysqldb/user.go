package mysqldb

import (
	"database/sql"
	"fmt"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/db"
)

type UserService struct {
	DB db.DB
}

func NewUserService(db db.DB) *UserService {
	return &UserService{DB: db}
}

// Ensure service implements interface.
var _ rlnotif.UserService = (*UserService)(nil)

func (s *UserService) User(id int64) (*rlnotif.User, error) {
	var user rlnotif.User

	row := s.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &user, fmt.Errorf("User %d: no such user", id)
		}
		return &user, fmt.Errorf("User %d: %v", id, err)
	}
	return &user, nil
}

func (s *UserService) Users() ([]*rlnotif.User, error) {
	var users []*rlnotif.User

	rows, err := s.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("Users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user rlnotif.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Users: %v", err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Users: %v", err)
	}
	return users, nil
}

func (s *UserService) CreateUser(user *rlnotif.User) (int64, error) {
	result, err := s.DB.Exec("INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}
	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %v", err)
	}

	return ID, nil
}
