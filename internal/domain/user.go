package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type UserRepository interface {
	CreateUser(user User) (string, error)
	LogIn(user User) (uuid.UUID, error)
}
