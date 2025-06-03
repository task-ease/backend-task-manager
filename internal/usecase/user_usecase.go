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

func (uc *UserUseCase) CreateUser(user domain.AuthUser) (string, error) {
	return uc.repo.CreateUser(user)
}

func (uc *UserUseCase) LogIn(user domain.AuthUser) (uuid.UUID, error) {
	return uc.repo.LogIn(user)
}

func (uc *UserUseCase) SearchUserByEmail(value string) ([]domain.User, error) {
	return uc.repo.SearchUserByEmail(value)
}
