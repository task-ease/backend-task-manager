package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) GetWorkSpaceTasks(workspaceId uuid.UUID) ([]response.GetAllWorkspaceTasks, error) {
	return uc.repo.GetWorkSpaceTasks(workspaceId)
}

func (uc *TaskUseCase) CreateTask(task request.CreateTask) (uuid.UUID, error) {
	return uc.repo.CreateTask(task)
}
