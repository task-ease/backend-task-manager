package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"time"
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

func (uc *TaskUseCase) UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error {
	return uc.repo.UpdateTaskColumn(taskId, columnId)
}

func (uc *TaskUseCase) UpdateTaskDescription(taskId uuid.UUID, description string) error {
	return uc.repo.UpdateTaskDescription(taskId, description)
}

func (uc *TaskUseCase) UpdateTaskIsFinished(taskId uuid.UUID, isFinished bool) error {
	return uc.repo.UpdateTaskIsFinished(taskId, isFinished)
}

func (uc *TaskUseCase) UpdateTaskDueDate(taskId uuid.UUID, dueDate time.Time) error {
	return uc.repo.UpdateTaskDueDate(taskId, dueDate)
}

func (uc *TaskUseCase) UpdateTaskPriority(taskId uuid.UUID, priority int) error {
	return uc.repo.UpdateTaskPriority(taskId, priority)
}
