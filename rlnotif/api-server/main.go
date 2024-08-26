package apiserver

import (
	"fmt"
	"time"

	"github.com/gchein/rate-limited-notification-service-go/rlnotif"
	"github.com/gchein/rate-limited-notification-service-go/rlnotif/mysqldb"
)

func Run() {
	u1 := &rlnotif.User{
		ID:        1,
		Name:      "susy",
		Email:     "susy@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	u2 := &rlnotif.User{
		ID:        2,
		Name:      "john",
		Email:     "john@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var repo mysqldb.UserRepo = []*rlnotif.User{u1, u2}
	// fmt.Printf("Value : %+v. Type: %T\n\n", repo, repo)

	s := mysqldb.NewUserService(&repo)
	// fmt.Printf("Value : %+v. Type: %T\n\n", s, s)
	// fmt.Printf("Value : %+v. Type: %T\n\n", s.DB, s.DB)

	for i := range *s.DB {
		// fmt.Printf("Value : %+v. Type: %T\n\n", v, v)

		u, _ := s.User(i + 1)
		fmt.Printf("Value : %+v. Type: %T\n\n", *u, *u)
	}

	fmt.Println("Adding a new user...")

	u3 := &rlnotif.User{
		ID:        3,
		Name:      "greg",
		Email:     "greg@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.CreateUser(u3)

	fmt.Println("Updated info:")

	for i := range *s.DB {
		// fmt.Printf("Value : %+v. Type: %T\n\n", v, v)

		u, _ := s.User(i + 1)
		fmt.Printf("Value : %+v. Type: %T\n\n", *u, *u)
	}

}
