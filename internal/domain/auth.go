package domain

import "github.com/google/uuid"

type AuthService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	VerifyToken(tokenString string) (string, error)
}
