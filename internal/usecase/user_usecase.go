package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
)

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) CreateUser(user domain.User) (string, error) {
	return uc.repo.CreateUser(user)
}

func (uc *UserUseCase) LogIn(user domain.User) (uuid.UUID, error) {
	return uc.repo.LogIn(user)
}
