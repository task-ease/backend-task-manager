package response

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
	"time"
)

type GetAllWorkspaceTasks struct {
	Id          uuid.UUID            `json:"id"`
	ColumnId    uuid.UUID            `json:"columnId"`
	ParentId    *uuid.UUID           `json:"parentId"`
	Type        enums.TaskTypes      `json:"type"`
	Title       string               `json:"title"`
	Description *string              `json:"description"`
	IsDone      bool                 `json:"isDone"`
	DueDate     *time.Time           `json:"dueDate"`
	Priority    enums.TaskPriorities `json:"priority"`
	Position    int                  `json:"position"`
	AssignedTo  []AssignedTaskUser   `json:"assignedTo"`
}

type AssignedTaskUser struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IconUrl  *string   `json:"iconUrl"`
}
