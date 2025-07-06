package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"time"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (r *taskRepository) CreateTask(task *domain.Task) (bool, error) {
	task.ID = uuid.New()
	task.CreatedAt = time.Now().UTC()
	task.UpdatedAt = time.Now().UTC()

	_, err := r.conn.Exec(context.Background(),
		`INSERT INTO tasks (
        id,
        workspace_id,
        author_id,
        created_at,
        title,
        description,
        is_finished,
        due_date,
        priority,
        updated_at,
		column_id
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		task.ID,
		task.WorkspaceID,
		task.AuthorID,
		task.CreatedAt,
		task.Title,
		task.Description,
		task.IsFinished,
		task.DueDate,
		task.Priority,
		task.UpdatedAt,
		task.ColumnID,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *taskRepository) CreateColumn(name string, workspaceId uuid.UUID, position int, color string) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		`INSERT INTO task_columns (
                          id, 
                          workspace_id,
                          name,
                          position,
                          color)
			VALUES ($1, $2, $3, $4, $5)`, uuid.New(), workspaceId, name, position, color)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *taskRepository) GetAllColumns(workspaceId uuid.UUID) ([]*domain.TaskColumn, error) {
	rows, err := r.conn.Query(context.Background(),
		`SELECT id, workspace_id, name, position, color FROM task_columns WHERE workspace_id = $1`, workspaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskColumns []*domain.TaskColumn
	for rows.Next() {
		var column domain.TaskColumn
		if err := rows.Scan(&column.ID, &column.WorkspaceId, &column.Name, &column.Position, &column.Color); err != nil {
			return nil, err
		}
		taskColumns = append(taskColumns, &column)
	}

	return taskColumns, nil
}

func (r *taskRepository) GetAllTasks(workspaceId uuid.UUID) ([]*domain.Task, error) {
	rows, err := r.conn.Query(context.Background(),
		`SELECT 
				id, 
				column_id, 
				author_id, 
				created_at, 
				title,
				description,
				is_finished,
				due_date,
				priority,
				updated_at
			FROM tasks WHERE workspace_id = $1`, workspaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskList []*domain.Task
	for rows.Next() {
		var (
			task        domain.Task
			description sql.NullString
			dueDate     sql.NullTime
			priority    sql.NullInt64
		)
		if err := rows.Scan(
			&task.ID,
			&task.ColumnID,
			&task.AuthorID,
			&task.CreatedAt,
			&task.Title,
			&description,
			&task.IsFinished,
			&dueDate,
			&priority,
			&task.UpdatedAt); err != nil {
			return nil, err
		}

		task.WorkspaceID = workspaceId

		if description.Valid {
			task.Description = &description.String
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		if priority.Valid {
			task.Priority = &priority.Int64
		}

		taskList = append(taskList, &task)
	}

	return taskList, nil
}

func (r *taskRepository) UpdateTaskTitle(taskId uuid.UUID, title string) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks SET title = $1, updated_at = NOW() WHERE id = $2`, title, taskId)
	return err
}

func (r *taskRepository) UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks SET column_id = $1, updated_at = NOW()  WHERE id = $2`, columnId, taskId)
	return err
}

func (r *taskRepository) UpdateTaskDescription(taskId uuid.UUID, description string) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks SET description = $1, updated_at = NOW() WHERE id = $2`, taskId, description)
	return err
}

func (r *taskRepository) UpdateTaskIsFinished(taskId uuid.UUID, isFinished bool) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks SET is_finished = $1, updated_at = NOW() WHERE id = $2`, isFinished, taskId)
	return err
}

func (r *taskRepository) UpdateTaskDueDate(taskId uuid.UUID, dueDate time.Time) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks SET due_date = $1, updated_at = NOW() WHERE id = $2`, dueDate, taskId)
	return err
}

func (r *taskRepository) UpdateTaskPriority(taskId uuid.UUID, priority int) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks SET priority = $1, updated_at = NOW() WHERE id = $2`, priority, taskId)
	return err
}
