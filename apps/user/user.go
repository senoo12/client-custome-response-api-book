package user

import (
	"time"
)

type LoginRequest struct {
	Email string
	Password string
}

type User struct {
	ID int
	Name string
	Email string
	Password string
	Created_At time.Time
	Updated_At time.Time
}

type UserResponse struct {
	ID int
	Name string
	Email string
	Created_At time.Time
	Updated_At time.Time
}

type LoginResponse struct {
	ID int
	Name string
	Email string
	Token string
}
