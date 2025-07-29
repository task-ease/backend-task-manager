package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/response"
)

type TaskRepository interface {
	GetWorkSpaceTasks(workspaceId uuid.UUID) ([]response.GetAllWorkspaceTasks, error)
}
