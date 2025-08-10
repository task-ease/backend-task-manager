package usecase

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/helper"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ColumnUsecase struct {
	columnRepo    domain.ColumnRepository
	baseRepo      *repository.BaseRepo
	workSpaceRepo domain.WorkSpaceRepository
}

func NewColumnUsecase(
	columnRepo domain.ColumnRepository,
	workSpaceRepo domain.WorkSpaceRepository,
	baseRepo *repository.BaseRepo,
) *ColumnUsecase {
	return &ColumnUsecase{columnRepo, baseRepo, workSpaceRepo}
}

func (uc *ColumnUsecase) GetColumns(ctx context.Context, workspaceId uuid.UUID, projectId, sprintId *uuid.UUID) ([]response.GetWorkspaceColumns, error) {
	return uc.columnRepo.GetColumns(ctx, workspaceId, projectId, sprintId)
}

func (uc *ColumnUsecase) CreateColumnTemplate(ctx context.Context, columnTmp request.CreateNewColumnTemplate) (uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (uuid.UUID, error) {
		return uc.CreateColumnTemplateTx(ctx, exec, columnTmp)
	})
}

func (uc *ColumnUsecase) CreateColumnTemplateTx(ctx context.Context, exec entities.Execer, columnTmp request.CreateNewColumnTemplate) (uuid.UUID, error) {
	id := uuid.New()

	if columnTmp.IsDone {
		if err := uc.columnRepo.ClearDoneFlagTx(ctx, exec, columnTmp.WorkspaceId); err != nil {
			return uuid.Nil, err
		}
	}

	if err := uc.columnRepo.CreateColumnTemplateTx(ctx, exec, columnTmp, id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (uc *ColumnUsecase) GetAllColumnTemplates(ctx context.Context, workSpaceId uuid.UUID) ([]entities.ColumnTemplate, error) {
	return uc.columnRepo.GetAllColumnTemplates(ctx, workSpaceId)
}

func (uc *ColumnUsecase) UpdateColumnTemplateName(ctx context.Context, columnTemplateId uuid.UUID, name string) error {
	return uc.columnRepo.UpdateColumnTemplateName(ctx, columnTemplateId, name)
}

func (uc *ColumnUsecase) UpdateColumnTemplateColor(ctx context.Context, columnTemplateId uuid.UUID, color string) error {
	return uc.columnRepo.UpdateColumnTemplateColor(ctx, columnTemplateId, color)
}

func (uc *ColumnUsecase) UpdateColumnTemplateStatusRequired(ctx context.Context, columnTemplateId uuid.UUID, status bool) error {
	return uc.columnRepo.UpdateColumnTemplateStatusRequired(ctx, columnTemplateId, status)
}

func (uc *ColumnUsecase) UpdateColumnTemplateStatusDone(ctx context.Context, columnTemplateId uuid.UUID) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		workspaceId, err := uc.workSpaceRepo.GetIdByColumnTemplateIdTx(ctx, exec, columnTemplateId)

		if err != nil {
			return err
		}

		if err = uc.columnRepo.ClearDoneFlagTx(ctx, exec, workspaceId); err != nil {
			return err
		}

		if err = uc.columnRepo.UpdateColumnTemplateSetDoneTx(ctx, exec, columnTemplateId); err != nil {
			return err
		}

		return uc.columnRepo.UpdateColumnTemplateStatusRequiredTx(ctx, exec, columnTemplateId, true)
	})
}

func (uc *ColumnUsecase) UpdateColumnTemplateStatusActive(ctx context.Context, columnTemplateId uuid.UUID, status bool) error {
	return uc.columnRepo.UpdateColumnTemplateStatusActive(ctx, columnTemplateId, status)
}

func (uc *ColumnUsecase) RenumberColumnTemplatesPositions(ctx context.Context, workspaceId uuid.UUID) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		ids, err := uc.columnRepo.GetColumnTemplatesIdsByWorkspaceOrderByPositionTx(ctx, exec, workspaceId)

		if err != nil {
			return err
		}

		return uc.columnRepo.RenumberColumnTemplatesPositionsTx(ctx, exec, ids)
	})
}
