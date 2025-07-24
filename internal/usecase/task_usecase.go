package usecase

import (
	"go-postgres-test/internal/domain"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}
