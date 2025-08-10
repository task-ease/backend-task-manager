package entities

import (
	"time"

	"github.com/google/uuid"
)

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
