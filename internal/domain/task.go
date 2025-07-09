package domain

import (
	"github.com/google/uuid"
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

type TaskColumn struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceId uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Position    int       `json:"position"`
	Color       string    `json:"color"`
	IsDone      bool      `json:"isDone"`
}

type TaskRepository interface {
	CreateTask(task *Task) (bool, error)
	CreateColumn(name string, workspaceId uuid.UUID, position int, color string) (bool, error)
	GetAllColumns(workspaceId uuid.UUID) ([]*TaskColumn, error)
	GetAllTasks(workspaceId uuid.UUID) ([]*Task, error)

	ReorderTasks(columnId uuid.UUID, orderedTaskIDs []uuid.UUID) error
	UpdateTask(task *Task) error
	MarkColumnAsDone(columnID uuid.UUID, isDone bool) error
	UpdateColumn(id, name, color string) error
}
