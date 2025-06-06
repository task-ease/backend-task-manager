package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID      `json:"id"`
	AuthorID    uuid.UUID      `json:"authorId"`
	ColumnID    uuid.UUID      `json:"columnId"`
	WorkspaceID uuid.UUID      `json:"workspaceId"`
	CreatedAt   time.Time      `json:"createdAt"`
	Title       string         `json:"title"`
	IsFinished  bool           `json:"isFinished"`
	Description sql.NullString `json:"description"`
	DueDate     sql.NullTime   `json:"dueDate"`
	Priority    sql.NullInt64  `json:"priority"`
	UpdatedAt   sql.NullTime   `json:"updatedAt"`
}

type TaskColumn struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Position    int       `json:"position"`
}

type TaskRepository interface {
	CreateTask(task *Task) (bool, error)
	CreateColumn(name string, workspaceId uuid.UUID, position int) (bool, error)
	GetAllColumns(workspaceId uuid.UUID) ([]*TaskColumn, error)
}
