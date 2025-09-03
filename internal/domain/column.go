package domain

import (
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/entities"
	"context"

	"github.com/google/uuid"
)

type ColumnRepository interface {
	GetColumns(ctx context.Context, workspaceId uuid.UUID, projectId, sprintId *uuid.UUID) ([]response.GetWorkspaceColumns, error)

	GetAllColumnTemplates(ctx context.Context, workSpaceId uuid.UUID) ([]entities.ColumnTemplate, error)
	UpdateColumnTemplateName(ctx context.Context, columnTemplateId uuid.UUID, name string) error
	UpdateColumnTemplateColor(ctx context.Context, columnTemplateId uuid.UUID, color string) error
	UpdateColumnTemplateStatusActive(ctx context.Context, columnId uuid.UUID, status bool) error
	UpdateColumnTemplateStatusRequired(ctx context.Context, columnTemplateId uuid.UUID, status bool) error

	AddColumnTx(ctx context.Context, exec entities.Execer, columnTemplateId, workspaceId uuid.UUID, projectId, sprintId *uuid.UUID) (uuid.UUID, error)
	RemoveColumnTx(ctx context.Context, exec entities.Execer, id uuid.UUID) error
	ClearDoneFlagTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID) error
	CreateColumnTemplateTx(ctx context.Context, exec entities.Execer, columnTmp request.CreateNewColumnTemplate, id uuid.UUID) error
	UpdateColumnTemplateSetDoneTx(ctx context.Context, exec entities.Execer, columnTemplateId uuid.UUID) error
	UpdateColumnTemplateStatusRequiredTx(ctx context.Context, exec entities.Execer, columnTemplateId uuid.UUID, status bool) error
	RenumberColumnTemplatesPositionsTx(ctx context.Context, exec entities.Execer, ids []uuid.UUID) error
	GetColumnTemplatesIdsByWorkspaceOrderByPositionTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID) ([]uuid.UUID, error)
}
