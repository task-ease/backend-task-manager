package usecase

import "go-postgres-test/internal/domain"

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) CreateUser(user domain.User) (string, error) {
	return uc.repo.CreateUser(user)
}
