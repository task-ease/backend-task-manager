package request

import (
	"backend-task-manager/internal/enums"
	"time"

	"github.com/google/uuid"
)

type CreateTask struct {
	ColumnId    uuid.UUID            `json:"columnId"`
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

type ChangeTaskPositionAndColumn struct {
	ToColumnId uuid.UUID  `json:"toColumnId"`
	PrevTaskId *uuid.UUID `json:"prevTaskId"`
}

type UpdateTaskType struct {
	TaskId uuid.UUID       `json:"taskId"`
	Value  enums.TaskTypes `json:"value"`
}
