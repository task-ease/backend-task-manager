package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
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

func (uc *ProjectUseCase) AddUserToProject(projectId uuid.UUID, userId uuid.UUID, role enums.ProjectRole) error {
	return uc.repo.AddUserToProject(projectId, userId, role)
}

func (uc *ProjectUseCase) GetAllUserProjects(userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error) {
	return uc.repo.GetAllUserProjects(userId, workspaceId)
}

func (uc *ProjectUseCase) GetUserRole(userId, projectId uuid.UUID) (enums.ProjectRole, error) {
	return uc.repo.GetUserRole(userId, projectId)
}

func (uc *ProjectUseCase) GetAllProjectMembers(projectId uuid.UUID) ([]response.GetAllProjectUsers, error) {
	return uc.repo.GetAllProjectMembers(projectId)
}

func (uc *ProjectUseCase) ChangeUserRole(userId, projectId uuid.UUID, role enums.ProjectRole) error {
	return uc.repo.ChangeUserRole(userId, projectId, role)
}
