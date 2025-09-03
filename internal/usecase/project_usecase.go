package usecase

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/dto"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/helper"
	"backend-task-manager/internal/repository"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProjectUseCase struct {
	projectRepo domain.ProjectRepository
	baseRepo    *repository.BaseRepo
	userRepo    domain.UserRepository
}

func NewProjectUseCase(projectRepo domain.ProjectRepository, baseRepo *repository.BaseRepo, userRepo domain.UserRepository) *ProjectUseCase {
	return &ProjectUseCase{projectRepo, baseRepo, userRepo}
}

func (uc *ProjectUseCase) CheckUserAccess(ctx context.Context, userId, projectId uuid.UUID) (dto.RolesMiddlewareDto, error) {
	role, err := uc.projectRepo.GetUserRole(ctx, userId, projectId)
	if err != nil {
		return dto.RolesMiddlewareDto{Role: enums.NoAccess, CanEdit: false}, err
	}

	return dto.RolesMiddlewareDto{Role: role, CanEdit: true}, nil
}

func (uc *ProjectUseCase) CreateProject(ctx context.Context, creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (uuid.UUID, error) {
		var id = uuid.New()
		if err := uc.projectRepo.CreateProjectTx(ctx, exec, creatorId, workSpaceId, id, name, prefix); err != nil {
			return uuid.Nil, err
		}

		if err := uc.projectRepo.AddUserToProjectTx(ctx, exec, id, creatorId, enums.ProjectCreator); err != nil {
			return uuid.Nil, err
		}
		return id, nil
	})
}

func (uc *ProjectUseCase) AddUserToProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID, role enums.UserRoles) (string, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (string, error) {
		if err := uc.projectRepo.AddUserToProject(ctx, projectId, userId, role); err != nil {
			return "", err
		}

		return uc.userRepo.GetEmailByUserId(ctx, userId)
	})
}

func (uc *ProjectUseCase) GetAllUserProjects(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error) {
	return uc.projectRepo.GetAllUserProjects(ctx, userId, workspaceId)
}

func (uc *ProjectUseCase) GetUserRole(ctx context.Context, userId, projectId uuid.UUID) (enums.UserRoles, error) {
	return uc.projectRepo.GetUserRole(ctx, userId, projectId)
}

func (uc *ProjectUseCase) GetAllProjectMembers(ctx context.Context, projectId uuid.UUID) ([]response.GetAllProjectUsers, error) {
	return uc.projectRepo.GetAllProjectMembers(ctx, projectId)
}

func (uc *ProjectUseCase) ChangeUserRole(ctx context.Context, userId, projectId uuid.UUID, role enums.UserRoles) error {
	return uc.projectRepo.ChangeUserRole(ctx, userId, projectId, role)
}

func (uc *ProjectUseCase) RemoveUserFromProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID) error {
	return uc.projectRepo.RemoveUserFromProject(ctx, projectId, userId)
}
