package domain

import (
	"context"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
	"time"

	"github.com/google/uuid"
)

type WorkSpace struct {
	ID        uuid.UUID `json:"id"`
	CreatorId uuid.UUID `json:"creator_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkSpaceRepository interface {
	GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]WorkSpace, error)
	GetWorkspaceName(ctx context.Context, workspaceId uuid.UUID) (string, error)

	AddUser(ctx context.Context, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error
	RemoveUser(ctx context.Context, workSpaceId, userId uuid.UUID) error
	GetAllMembers(ctx context.Context, workSpaceId uuid.UUID) ([]entities.MemberUser, error)
	ChangeUserRole(ctx context.Context, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error
	SearchWorkspaceMember(ctx context.Context, workSpaceId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error)

	AddUserTx(ctx context.Context, exec entities.Execer, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error
	CreateWorkSpaceTx(ctx context.Context, exec entities.Execer, workSpace WorkSpace) error
	HasUserWorkspaceTx(ctx context.Context, exec entities.Execer, userId, workspaceId uuid.UUID, exists *bool) error
	GetUserRoleTx(ctx context.Context, exec entities.Execer, userId, workspaceId uuid.UUID) (enums.WorkspaceRole, error)
	GetIdByColumnTemplateIdTx(ctx context.Context, exec entities.Execer, columnTemplateId uuid.UUID) (uuid.UUID, error)
}
