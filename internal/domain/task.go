package domain

import (
	"context"
	"github.com/google/uuid"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"
)

type TaskRepository interface {
	GetWorkSpaceTasks(workspaceId uuid.UUID) ([]response.GetAllWorkspaceTasks, error)
	CreateTask(task request.CreateTask) (uuid.UUID, error)

	UpdateTaskTitle(taskId uuid.UUID, value string) error
	UpdateTaskDescription(taskId uuid.UUID, value string) error
	UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error
	UpdateTaskPriority(taskId uuid.UUID, priority enums.TaskPriorities) error
	UpdateTaskAssigned(taskId uuid.UUID, userId uuid.UUID) error

	UpdateTaskColumnTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error
	UpdateTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error
}
