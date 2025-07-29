package entities

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
	"time"
)

type Task struct {
	Id          uuid.UUID            `json:"id"`
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
	DeletedAt   *time.Time           `json:"deletedAt"`
	DueDate     *time.Time           `json:"dueDate"`
	Priority    enums.TaskPriorities `json:"priority"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Position    int                  `json:"position"`
}

type Column struct {
	Id         uuid.UUID `json:"id"`
	ProjectId  uuid.UUID `json:"projectId"`
	TemplateId uuid.UUID `json:"templateId"`
}

type ColumnTemplate struct {
	Id          uuid.UUID `json:"id"`
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Position    int       `json:"position"`
	IsRequired  bool      `json:"isRequired"`
	IsActive    bool      `json:"isActive"`
	IsDone      bool      `json:"isDone"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	GlobalTasks bool      `json:"globalTasks"`
}
