package domain

import (
	"context"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user entities.AuthUser) (string, error)
	LogIn(user entities.AuthUser) (uuid.UUID, error)
	SearchUserByEmail(value string) ([]entities.User, error)
	ChangeOnlineStatus(userId uuid.UUID, status bool) error
	GetWorkspaceUserRole(userID uuid.UUID, workspaceId uuid.UUID) (enums.WorkspaceRole, error)

	CreateUserTx(ctx context.Context, exec entities.Execer, user entities.AuthUser) (string, error)
}
