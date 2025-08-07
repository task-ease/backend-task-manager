package usecase

import (
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/request/query"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) GetAllTasks(data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error) {
	return uc.repo.GetALlTasks(data)
}

func (uc *TaskUseCase) CreateTask(task request.CreateTask) (response.CreateTask, error) {
	return uc.repo.CreateTask(task)
}

func (uc *TaskUseCase) UpdateTaskTitle(taskId uuid.UUID, value string) error {
	return uc.repo.UpdateTaskTitle(taskId, value)
}

func (uc *TaskUseCase) UpdateTaskDescription(taskId uuid.UUID, value string) error {
	return uc.repo.UpdateTaskDescription(taskId, value)
}

func (uc *TaskUseCase) UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error {
	return uc.repo.UpdateTaskColumn(taskId, columnId)
}

func (uc *TaskUseCase) UpdateTaskPriority(taskId uuid.UUID, priority enums.TaskPriorities) error {
	return uc.repo.UpdateTaskPriority(taskId, priority)
}

func (uc *TaskUseCase) UpdateTaskAssigned(taskId uuid.UUID, userId uuid.UUID) error {
	return uc.repo.UpdateTaskAssigned(taskId, userId)
}

func (uc *TaskUseCase) RemoveTaskAssigned(taskId uuid.UUID) error {
	return uc.repo.RemoveTaskAssigned(taskId)
}
