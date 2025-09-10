package domain

import (
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/request/query"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/entities"
	"backend-task-manager/internal/enums"
	"context"

	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateNew(ctx context.Context, task request.CreateTask, id uuid.UUID) (response.CreateTask, error)
	GetAll(ctx context.Context, data query.TaskLocationWithSprintQuery) ([]response.GetAllWorkspaceTasks, error)
	Search(ctx context.Context, exec entities.Execer, data query.TaskLocationQuery, value string) ([]response.SearchTasks, error)

	UpdateType(ctx context.Context, taskId uuid.UUID, value enums.TaskTypes) error
	UpdateTitle(ctx context.Context, taskId uuid.UUID, value string) error
	UpdatePriority(ctx context.Context, taskId uuid.UUID, priority enums.TaskPriorities) error
	RemoveAssigned(ctx context.Context, taskId uuid.UUID) error
	UpdateDescription(ctx context.Context, taskId uuid.UUID, value string) error

	IsDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error
	GetAllTx(ctx context.Context, exec entities.Execer, data query.TaskLocationWithSprintQuery) ([]response.GetAllWorkspaceTasks, error)
	GetTypeTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) (enums.TaskTypes, error)
	GetPrefixTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID, projectId *uuid.UUID, prefix *string) (err error)
	CreateNewTx(ctx context.Context, exec entities.Execer, task request.CreateTask, id uuid.UUID) (response.CreateTask, error)
	UpdateTypeTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, value enums.TaskTypes) error
	HasAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, has *bool) error
	GetLocationTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) (query.TaskLocationQuery, error)
	IsColumnDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error
	UpdateParentIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, parentId *uuid.UUID) error
	IfExistsParentTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) (bool, error)
	RemoveAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) error
	GetPrefixNumberTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, prefixNumber *int) error
	UpdateChildrenIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, parentId *uuid.UUID) error
	ChangePositionAndColumnTx(ctx context.Context, exec entities.Execer, dto request.ChangeTaskPositionAndColumn) (float64, error)

	UpdateColumnIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error
	UpdateAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error
	UpdateDoneStatusTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, to bool) error
}
