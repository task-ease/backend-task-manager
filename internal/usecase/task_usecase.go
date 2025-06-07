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

func (uc *TaskUseCase) GetAllTasks(workspaceId uuid.UUID) ([]*domain.Task, error) {
	return uc.repo.GetAllTasks(workspaceId)
}

func (uc *TaskUseCase) UpdateTaskTitle(taskId uuid.UUID, title string) error {
	return uc.repo.UpdateTaskTitle(taskId, title)
}
