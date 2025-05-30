package usecase

import "go-postgres-test/internal/domain"

type AuthUseCase struct {
	auth domain.AuthService
}

func NewAuthUseCase(auth domain.AuthService) *AuthUseCase {
	return &AuthUseCase{auth: auth}
}
