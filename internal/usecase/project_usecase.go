package usecase

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/helper"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProjectUseCase struct {
	repo     domain.ProjectRepository
	baseRepo *repository.BaseRepo
}

func NewProjectUseCase(repo domain.ProjectRepository, baseRepo *repository.BaseRepo) *ProjectUseCase {
	return &ProjectUseCase{repo, baseRepo}
}

func (uc *ProjectUseCase) CreateProject(ctx context.Context, creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (uuid.UUID, error) {
		var id = uuid.New()
		if err := uc.repo.CreateProjectTx(ctx, exec, creatorId, workSpaceId, id, name, prefix); err != nil {
			return uuid.Nil, err
		}

		if err := uc.repo.AddUserToProjectTx(ctx, exec, id, creatorId, enums.ProjectRoleCreator); err != nil {
			return uuid.Nil, err
		}
		return id, nil
	})
}

func (uc *ProjectUseCase) AddUserToProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID, role enums.ProjectRole) error {
	return uc.repo.AddUserToProject(ctx, projectId, userId, role)
}

func (uc *ProjectUseCase) GetAllUserProjects(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error) {
	return uc.repo.GetAllUserProjects(ctx, userId, workspaceId)
}

func (uc *ProjectUseCase) GetUserRole(ctx context.Context, userId, projectId uuid.UUID) (enums.ProjectRole, error) {
	return uc.repo.GetUserRole(ctx, userId, projectId)
}

func (uc *ProjectUseCase) GetAllProjectMembers(ctx context.Context, projectId uuid.UUID) ([]response.GetAllProjectUsers, error) {
	return uc.repo.GetAllProjectMembers(ctx, projectId)
}

func (uc *ProjectUseCase) ChangeUserRole(ctx context.Context, userId, projectId uuid.UUID, role enums.ProjectRole) error {
	return uc.repo.ChangeUserRole(ctx, userId, projectId, role)
}
