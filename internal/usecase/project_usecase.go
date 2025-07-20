package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types/user"
)

type ProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewProjectUseCase(repo domain.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{repo}
}

func (uc *ProjectUseCase) CreateProject(creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error) {
	return uc.repo.CreateProject(creatorId, workSpaceId, name, prefix)
}

func (uc *ProjectUseCase) AddUserToProject(projectId uuid.UUID, userId uuid.UUID, role user.ProjectRole) error {
	return uc.repo.AddUserToProject(projectId, userId, role)
}

func (uc *ProjectUseCase) GetAllUserProjects(userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error) {
	return uc.repo.GetAllUserProjects(userId, workspaceId)
}

func (uc *ProjectUseCase) GetUserRole(userId, projectId uuid.UUID) (user.ProjectRole, error) {
	return uc.repo.GetUserRole(userId, projectId)
}
