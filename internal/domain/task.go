package domain

import (
	"context"
	"go-postgres-test/internal/dto/request"
	"go-postgres-test/internal/dto/request/query"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task request.CreateTask, id uuid.UUID) (response.CreateTask, error)
	GetALlTasks(ctx context.Context, data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error)

	UpdateTaskTitle(ctx context.Context, taskId uuid.UUID, value string) error
	UpdateTaskPriority(ctx context.Context, taskId uuid.UUID, priority enums.TaskPriorities) error
	RemoveTaskAssigned(ctx context.Context, taskId uuid.UUID) error
	UpdateTaskDescription(ctx context.Context, taskId uuid.UUID, value string) error
	UpdateTaskType(ctx context.Context, taskId uuid.UUID, value enums.TaskTypes) error

	GetPrefixTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID, projectId *uuid.UUID, prefix *string) (err error)
	IsTaskDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error
	CreateTaskTx(ctx context.Context, exec entities.Execer, task request.CreateTask, id uuid.UUID) (response.CreateTask, error)
	GetALlTasksTx(ctx context.Context, exec entities.Execer, data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error)
	IsColumnDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error
	GetPrefixNumberTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, prefixNumber *int) error
	IsTaskHasAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, has *bool) error
	RemoveTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) error
	ChangeTaskPositionAndColumnTx(ctx context.Context, exec entities.Execer, dto request.ChangeTaskPositionAndColumn) (float64, error)

	UpdateTaskColumnIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error
	UpdateTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error
	UpdateTaskDoneStatusTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, to bool) error
}
