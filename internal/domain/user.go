package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
type AuthUser struct {
	User
	Password string `json:"password"`
}
type MemberUser struct {
	User
	JoinedAt time.Time `json:"joined_at"`
	Role     string    `json:"role"`
	Position *string   `json:"position"`
}

type UserRepository interface {
	CreateUser(user AuthUser) (string, error)
	LogIn(user AuthUser) (uuid.UUID, error)
	SearchUserByEmail(value string) ([]User, error)
}
