package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase { return &TaskUseCase{repo: repo} }

func (uc *TaskUseCase) CreateTask(task *domain.Task) (bool, error) { return uc.repo.CreateTask(task) }

func (uc *TaskUseCase) GetAllColumns(workspaceId uuid.UUID) ([]*domain.TaskColumn, error) {
	return uc.repo.GetAllColumns(workspaceId)
}
