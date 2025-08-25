package entities

import (
	"go-postgres-test/internal/enums"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id           uuid.UUID            `json:"id"`
	ColumnId     uuid.UUID            `json:"columnId"`
	WorkspaceId  uuid.UUID            `json:"workspaceId"`
	ProjectId    *uuid.UUID           `json:"projectId"`
	SprintId     *uuid.UUID           `json:"sprintId"`
	AuthorId     uuid.UUID            `json:"authorId"`
	ParentId     *uuid.UUID           `json:"parentId"`
	Type         enums.TaskTypes      `json:"type"`
	Title        string               `json:"title"`
	Description  *string              `json:"description"`
	IsDone       bool                 `json:"isDone"`
	DeletedAt    *time.Time           `json:"deletedAt"`
	DueDate      *time.Time           `json:"dueDate"`
	Priority     enums.TaskPriorities `json:"priority"`
	CreatedAt    time.Time            `json:"createdAt"`
	UpdatedAt    time.Time            `json:"updatedAt"`
	Position     float64              `json:"position"`
	PrefixNumber *int                 `json:"prefixNumber"`
}
