package request

import (
	"backend-task-manager/internal/enums"

	"github.com/google/uuid"
)

type CreateNewColumnTemplate struct {
	Color       string    `json:"color"`
	Name        string    `json:"name"`
	IsRequired  bool      `json:"isRequired"`
	IsDone      bool      `json:"isDone"`
	Position    int       `json:"position"`
	WorkspaceId uuid.UUID `json:"workspaceId"`
	GlobalTasks bool      `json:"globalTasks"`
}

type UpdateColumnTemplate struct {
	Value            string                          `json:"value"`
	ColumnTemplateId uuid.UUID                       `json:"columnTemplateId"`
	Method           enums.UpdateColumnTemplateValue `json:"method"`
}

type UpdateColumnTemplateStatus struct {
	Value            bool                             `json:"Value"`
	ColumnTemplateId uuid.UUID                        `json:"columnTemplateId"`
	Method           enums.UpdateColumnTemplateStatus `json:"method"`
}
