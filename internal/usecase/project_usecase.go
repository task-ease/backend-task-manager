package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
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
