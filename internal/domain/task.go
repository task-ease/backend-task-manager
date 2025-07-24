package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
	"time"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	AuthorID    uuid.UUID  `json:"authorId"`
	ColumnID    uuid.UUID  `json:"columnId"`
	WorkspaceID uuid.UUID  `json:"workspaceId"`
	CreatedAt   time.Time  `json:"createdAt"`
	Title       string     `json:"title"`
	IsFinished  bool       `json:"isFinished"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	Status      *int       `json:"status"`
	Priority    *int       `json:"priority"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Position    int        `json:"position"`
}

type TaskRepository interface {
	CreateColumnTemplate(workspaceId uuid.UUID, columnTmp enums.ColumnTemplate) error
}
