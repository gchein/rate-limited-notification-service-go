package mysqldb

import (
	"database/sql"
	"fmt"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) User(id int) (*rlnotif.User, error) {
	db := s.DB

	var user rlnotif.User

	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
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
	db := s.DB
	var users []*rlnotif.User

	rows, err := db.Query("SELECT * FROM users")
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

func (s *UserService) CreateUser(user *rlnotif.User) error {
	db := s.DB

	result, err := db.Exec("INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}

	return nil
}
