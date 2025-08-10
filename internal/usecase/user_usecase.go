package usecase

import (
	"context"
	"fmt"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/helper"
	"go-postgres-test/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo domain.UserRepository
	baseRepo *repository.BaseRepo
}

func NewUserUsecase(r domain.UserRepository, baseRepo *repository.BaseRepo) *UserUseCase {
	return &UserUseCase{userRepo: r, baseRepo: baseRepo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user entities.AuthUser) (uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (uuid.UUID, error) {
		var exists bool

		if err := uc.userRepo.CheckIfExistsByEmailTx(ctx, exec, user.Email, &exists); err != nil {
			return uuid.Nil, err
		}

		if exists {
			return uuid.Nil, fmt.Errorf("email already exists")
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			return uuid.Nil, err
		}

		id := uuid.New()

		if err := uc.userRepo.CreateUserTx(ctx, exec, user, id, passwordHash); err != nil {
			return uuid.Nil, err
		}

		return id, nil
	})
}

func (uc *UserUseCase) LogIn(ctx context.Context, user entities.AuthUser) (uuid.UUID, error) {
	dto, err := uc.userRepo.GetIdAndPasswordHash(ctx, user.Email)
	if err != nil {
		return uuid.Nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dto.PasswordHash), []byte(user.Password))

	if err != nil {
		return uuid.Nil, err
	}

	return dto.ID, nil
}

func (uc *UserUseCase) SearchUserByEmail(ctx context.Context, value string) ([]entities.User, error) {
	return uc.userRepo.SearchUserByEmail(ctx, value)
}

func (uc *UserUseCase) GetWorkspaceUserRole(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (enums.WorkspaceRole, error) {
	return uc.userRepo.GetWorkspaceUserRole(ctx, userID, workspaceID)
}
