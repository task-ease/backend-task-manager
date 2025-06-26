package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) CreateTask(task *domain.Task) (bool, error) {
	return uc.repo.CreateTask(task)
}

func (uc *TaskUseCase) GetAllColumns(workspaceId uuid.UUID) ([]*domain.TaskColumn, error) {
	return uc.repo.GetAllColumns(workspaceId)
}

func (uc *TaskUseCase) GetAllTasks(workspaceId uuid.UUID) ([]*domain.Task, error) {
	return uc.repo.GetAllTasks(workspaceId)
}

func (uc *TaskUseCase) UpdateTask(task *domain.Task) error {
	return uc.repo.UpdateTask(task)
}

func (uc *TaskUseCase) ReorderTasks(columnId uuid.UUID, orderedTaskIDs []uuid.UUID) error {
	return uc.repo.ReorderTasks(columnId, orderedTaskIDs)
}

func (uc *TaskUseCase) MarkColumnAsDone(columnID uuid.UUID, isDone bool) error {
	return uc.repo.MarkColumnAsDone(columnID, isDone)
}

func (uc *TaskUseCase) UpdateColumn(id, name, color string) error {
	return uc.repo.UpdateColumn(id, name, color)
}
