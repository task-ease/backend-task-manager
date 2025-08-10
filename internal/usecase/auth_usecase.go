package usecase

import "go-postgres-test/internal/domain"

type AuthUseCase struct {
	auth domain.AuthService
}
