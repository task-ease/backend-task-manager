package domain

import (
	"context"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	GetUserRole(ctx context.Context, userId, projectId uuid.UUID) (enums.UserRoles, error)
	CreateProject(ctx context.Context, creatorId, workSpaceId, projectId uuid.UUID, name, prefix string) error
	ChangeUserRole(ctx context.Context, userId, projectId uuid.UUID, role enums.UserRoles) error
	AddUserToProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID, role enums.UserRoles) error
	GetAllUserProjects(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error)
	GetAllProjectMembers(ctx context.Context, projectId uuid.UUID) ([]response.GetAllProjectUsers, error)
	RemoveUserFromProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID) error

	GetUserRoleTx(ctx context.Context, exec entities.Execer, userId, projectId uuid.UUID) (enums.UserRoles, error)
	CreateProjectTx(ctx context.Context, exec entities.Execer, creatorId, workSpaceId, projectId uuid.UUID, name, prefix string) error
	AddUserToProjectTx(ctx context.Context, exec entities.Execer, projectId uuid.UUID, userId uuid.UUID, role enums.UserRoles) error
	GetIdByDocumentIdTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID) (uuid.UUID, error)
}
