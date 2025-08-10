package domain

import (
	"context"
	"go-postgres-test/internal/dto"
	"go-postgres-test/internal/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entities.AuthUser, id uuid.UUID, passwordHash []byte) error
	SearchUserByEmail(ctx context.Context, value string) ([]entities.User, error)
	ChangeOnlineStatus(ctx context.Context, userId uuid.UUID, status bool) error
	CheckIfExistsByEmail(ctx context.Context, email string, exists *bool) error
	GetIdAndPasswordHash(ctx context.Context, email string) (dto.UserIdAndPasswordHash, error)

	GetIdAndPasswordHashTx(ctx context.Context, exec entities.Execer, email string) (dto.UserIdAndPasswordHash, error)
	CheckIfExistsByEmailTx(ctx context.Context, exec entities.Execer, email string, exists *bool) error
	CreateUserTx(ctx context.Context, exec entities.Execer, user entities.AuthUser, id uuid.UUID, passwordHash []byte) error
}
