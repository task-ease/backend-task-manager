package entities

import (
	"go-postgres-test/internal/enums"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	UserIconUrl  *string   `json:"userIconUrl"`
	IsOnline     bool      `json:"isOnline"`
	LastOnlineAt time.Time `json:"lastOnlineAt"`
}

type AuthUser struct {
	User
	Password string `json:"password"`
}

type MemberUser struct {
	User
	JoinedAt time.Time       `json:"joined_at"`
	Role     enums.UserRoles `json:"role"`
	Position *string         `json:"position"`
}
