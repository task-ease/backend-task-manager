package request

import (
	"github.com/google/uuid"
)

type UpdateColumnStatusRequest struct {
	Status   bool      `json:"status"`
	ColumnId uuid.UUID `json:"columnId"`
}

type UpdateColumnNameRequest struct {
	ColumnId uuid.UUID `json:"columnId"`
	Name     string    `json:"name"`
}

type CreateNewColumnTemplate struct {
	Color       string    `json:"color"`
	Name        string    `json:"name"`
	IsRequired  bool      `json:"isRequired"`
	IsDone      bool      `json:"isDone"`
	Position    int       `json:"position"`
	WorkspaceId uuid.UUID `json:"workspaceId"`
	GlobalTasks bool      `json:"globalTasks"`
}
