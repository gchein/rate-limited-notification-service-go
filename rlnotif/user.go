package rlnotif

import (
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserService interface {
	User(id int) (*User, error)
	Users() ([]*User, error)
	CreateUser(user *User) error
}
