package usecase

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/helper"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WorkSpaceUsecase struct {
	workspaceRepo domain.WorkSpaceRepository
	columnUsecase *ColumnUsecase
	baseRepo      *repository.BaseRepo
}

func NewWorkSpaceUsecase(workspaceRepo domain.WorkSpaceRepository, columnUsecase *ColumnUsecase, baseRepo *repository.BaseRepo) *WorkSpaceUsecase {
	return &WorkSpaceUsecase{workspaceRepo, columnUsecase, baseRepo}
}

func (uc *WorkSpaceUsecase) CreateWorkSpace(ctx context.Context, workspace domain.WorkSpace) (uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (uuid.UUID, error) {
		workspace.ID = uuid.New()

		if err := uc.workspaceRepo.CreateWorkSpaceTx(ctx, exec, workspace); err != nil {
			return uuid.Nil, err
		}

		if err := uc.workspaceRepo.AddUserTx(ctx, exec, workspace.ID, workspace.CreatorId, enums.WorkspaceCreator); err != nil {
			return uuid.Nil, err
		}

		columnList := []request.CreateNewColumnTemplate{
			{
				WorkspaceId: workspace.ID,
				Name:        "To Do",
				Color:       "#787878",
				Position:    0,
				IsRequired:  true,
				IsDone:      false,
				GlobalTasks: true,
			},
			{
				WorkspaceId: workspace.ID,
				Name:        "In Progress",
				Color:       "#3d66b8",
				Position:    10,
				IsRequired:  true,
				IsDone:      false,
				GlobalTasks: true,
			},
			{
				WorkspaceId: workspace.ID,
				Name:        "Done",
				Color:       "#00BFA6",
				Position:    20,
				IsRequired:  true,
				IsDone:      true,
				GlobalTasks: true,
			},
		}

		for _, column := range columnList {
			_, err := uc.columnUsecase.CreateColumnTemplateTx(ctx, exec, column)
			if err != nil {
				return uuid.Nil, err
			}
		}

		return workspace.ID, nil
	})
}

func (uc *WorkSpaceUsecase) GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.WorkSpace, error) {
	return uc.workspaceRepo.GetAllByUserId(ctx, userId)
}

func (uc *WorkSpaceUsecase) AddUser(ctx context.Context, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error {
	return uc.workspaceRepo.AddUser(ctx, workSpaceId, userId, role)
}

func (uc *WorkSpaceUsecase) GetAllMembers(ctx context.Context, workSpaceId uuid.UUID) ([]entities.MemberUser, error) {
	return uc.workspaceRepo.GetAllMembers(ctx, workSpaceId)
}

func (uc *WorkSpaceUsecase) RemoveUser(ctx context.Context, workSpaceId, userId uuid.UUID) error {
	return uc.workspaceRepo.RemoveUser(ctx, workSpaceId, userId)
}

func (uc *WorkSpaceUsecase) HasUserWorkspace(ctx context.Context, userId, workspaceId uuid.UUID) (enums.WorkspaceRole, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (enums.WorkspaceRole, error) {
		var exists bool
		if err := uc.workspaceRepo.HasUserWorkspaceTx(ctx, exec, userId, workspaceId, &exists); err != nil {
			return enums.WorkspaceNotAllowed, err
		}

		return uc.workspaceRepo.GetUserRoleTx(ctx, exec, userId, workspaceId)
	})
}

func (uc *WorkSpaceUsecase) ChangeUserRole(ctx context.Context, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error {
	return uc.workspaceRepo.ChangeUserRole(ctx, workSpaceId, userId, role)
}

func (uc *WorkSpaceUsecase) SearchWorkspaceMember(ctx context.Context, workSpaceId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error) {
	return uc.workspaceRepo.SearchWorkspaceMember(ctx, workSpaceId, value)
}

func (uc *WorkSpaceUsecase) GetWorkspaceName(ctx context.Context, workspaceId uuid.UUID) (string, error) {
	return uc.workspaceRepo.GetWorkspaceName(ctx, workspaceId)
}
