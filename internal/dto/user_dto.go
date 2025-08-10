package dto

import "github.com/google/uuid"

type UserIdAndPasswordHash struct {
	ID           uuid.UUID `json:"id"`
	PasswordHash string    `json:"passwordHash"`
}
