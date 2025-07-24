package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) CreateColumnTemplate(columnTmp entities.ColumnTemplate) (uuid.UUID, error) {
	return uc.repo.CreateColumnTemplate(columnTmp)
}
