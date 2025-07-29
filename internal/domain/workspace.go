package domain

import (
	"context"
	"github.com/google/uuid"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
	"time"
)

type WorkSpace struct {
	ID        uuid.UUID `json:"id"`
	CreatorId uuid.UUID `json:"creator_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkSpaceRepository interface {
	CreateWorkSpace(workspace WorkSpace) (uuid.UUID, error)
	GetAllUserSpaces(userId uuid.UUID) ([]WorkSpace, error)
	GetWorkSpaceColumns(workspaceId uuid.UUID) ([]response.GetWorkSpaceColumns, error)

	AddUserToWorkSpace(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error)
	GetAllSpaceMembers(workSpaceId uuid.UUID) ([]entities.MemberUser, error)
	RemoveUser(workSpaceId string, userId string) (bool, error)
	HasUserWorkspace(userId string, workspaceId string) (enums.WorkspaceRole, error)
	ChangeUserRole(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error)
	SearchWorkspaceMember(workSpaceId, userId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error)

	CreateColumnTemplate(columnTmp entities.ColumnTemplate) (uuid.UUID, error)
	GetAllColumnTemplates(workSpaceId uuid.UUID) ([]entities.ColumnTemplate, error)
	UpdateColumnTemplateStatusRequired(columnId uuid.UUID, status bool) error
	UpdateColumnTemplateStatusDone(columnId uuid.UUID) error
	UpdateColumnTemplateStatusActive(columnId uuid.UUID, status bool) error
	UpdateColumnTemplateStatusGlobalTasks(columnId uuid.UUID, status bool) error
	UpdateColumnTemplateName(columnId uuid.UUID, name string) error
	RenumberColumnTemplatesPositions(workspaceId uuid.UUID) error

	CreateWorkSpaceTx(ctx context.Context, exec entities.Execer, workSpace WorkSpace) (uuid.UUID, error)
	AddUserToWorkSpaceTx(ctx context.Context, exec entities.Execer, workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error)
	HasUserWorkspaceTx(ctx context.Context, exec entities.Execer, userId string, workspaceId string) (enums.WorkspaceRole, error)
	CreateColumnTemplateTx(ctx context.Context, exec entities.Execer, columnTmp entities.ColumnTemplate) (uuid.UUID, error)
	UpdateColumnTemplateStatusDoneTx(ctx context.Context, exec entities.Execer, columnId uuid.UUID) error
	RenumberColumnTemplatesPositionsTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID) error
}
