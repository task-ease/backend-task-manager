package request

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
	"time"
)

type CreateTask struct {
	ColumnId    uuid.UUID            `json:"columnId"`
	WorkspaceId uuid.UUID            `json:"workspaceId"`
	ProjectId   *uuid.UUID           `json:"projectId"`
	SprintId    *uuid.UUID           `json:"sprintId"`
	AuthorId    uuid.UUID            `json:"authorId"`
	ParentId    *uuid.UUID           `json:"parentId"`
	Type        enums.TaskTypes      `json:"type"`
	Title       string               `json:"title"`
	Description *string              `json:"description"`
	IsDone      bool                 `json:"isDone"`
	DueDate     *time.Time           `json:"dueDate"`
	Priority    enums.TaskPriorities `json:"priority"`
	Position    int                  `json:"position"`
}
