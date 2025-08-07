package domain

import (
	"context"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/request/query"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
)

type TaskRepository interface {
	GetALlTasks(data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error)
	CreateTask(task request.CreateTask) (response.CreateTask, error)

	UpdateTaskTitle(taskId uuid.UUID, value string) error
	UpdateTaskDescription(taskId uuid.UUID, value string) error
	UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error
	UpdateTaskPriority(taskId uuid.UUID, priority enums.TaskPriorities) error
	UpdateTaskAssigned(taskId uuid.UUID, userId uuid.UUID) error

	GetALlTasksTx(ctx context.Context, exec entities.Execer, data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error)
	CreateTaskTx(ctx context.Context, exec entities.Execer, task request.CreateTask) (response.CreateTask, error)

	UpdateTaskColumnTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error
	UpdateTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error
}
