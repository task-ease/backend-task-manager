package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"log"
	"time"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (r *taskRepository) CreateTask(task *domain.Task) (bool, error) {
	task.ID = uuid.New()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	var maxPos int
	err := r.conn.QueryRow(
		context.Background(),
		`SELECT MAX(position) FROM tasks WHERE column_id = $1`,
		task.ColumnID,
	).Scan(&maxPos)
	if err != nil {
		return false, err
	}

	_, err = r.conn.Exec(context.Background(),
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
    status,
    updated_at,
	column_id,
    position
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		task.ID,
		task.WorkspaceID,
		task.AuthorID,
		task.CreatedAt,
		task.Title,
		task.Description,
		task.IsFinished,
		task.DueDate,
		task.Priority,
		task.Status,
		task.UpdatedAt,
		task.ColumnID,
		task.Position,
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
	status,
	position,
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
			status      sql.NullInt64
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
			&status,
			&task.Position,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if description.Valid {
			task.Description = &description.String
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		if status.Valid {
			val := int(status.Int64)
			task.Status = &val
		}

		taskList = append(taskList, &task)
	}

	return taskList, nil
}

func (r *taskRepository) UpdateTaskTitle(taskId uuid.UUID, title string) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             title = $1,
             updated_at = NOW()
         WHERE id = $2`,
		title,
		taskId,
	)
	return err
}

func (r *taskRepository) UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             column_id = $1,
             updated_at = NOW()
         WHERE id = $2`,
		columnId,
		taskId,
	)
	return err
}

func (r *taskRepository) UpdateTaskDescription(taskId uuid.UUID, description string) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             description = $1,
             updated_at = NOW()
         WHERE id = $2`,
		description,
		taskId,
	)
	return err
}

func (r *taskRepository) UpdateTaskIsFinished(taskId uuid.UUID, isFinished bool) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             is_finished = $1,
             updated_at = NOW()
         WHERE id = $2`,
		isFinished,
		taskId,
	)
	return err
}

func (r *taskRepository) UpdateTaskDueDate(taskId uuid.UUID, dueDate time.Time) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             due_date = $1,
             updated_at = NOW()
         WHERE id = $2`,
		dueDate,
		taskId,
	)
	return err
}

func (r *taskRepository) UpdateTaskPriority(taskId uuid.UUID, priority int) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             priority = $1,
             updated_at = NOW()
         WHERE id = $2`,
		priority,
		taskId,
	)
	return err
}
func (r *taskRepository) UpdateTask(task *domain.Task) error {
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		`UPDATE tasks
         SET position = new_position
         FROM (
             SELECT id, ROW_NUMBER() OVER (ORDER BY position) - 1 AS new_position
             FROM tasks
             WHERE column_id = $1
         ) AS ranked
         WHERE tasks.id = ranked.id`,
		task.ColumnID,
	)
	if err != nil {
		log.Printf("UpdateTask error: %v", err)
		return err
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE tasks 
         SET 
             title = $1,
             description = $2,
             is_finished = $3,
             due_date = $4,
             priority = $5,
             status = $6,
             column_id = $7,
             updated_at = NOW()
         WHERE id = $8`,
		task.Title,
		task.Description,
		task.IsFinished,
		task.DueDate,
		task.Priority,
		task.Status,
		task.ColumnID,
		task.ID,
	)
	if err != nil {
		log.Printf("UpdateTask error: %v", err)
		return err
	}

	return tx.Commit(context.Background())
}

func (r *taskRepository) UpdateTaskPosition(taskId uuid.UUID, position int) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE tasks 
         SET position = $1,
             updated_at = NOW()
         WHERE id = $2`,
		position,
		taskId,
	)
	return err
}

func (r *taskRepository) ReorderTasks(columnId uuid.UUID, orderedTaskIDs []uuid.UUID) error {
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	log.Printf("ReorderTasks: columnId=%s, orderedTaskIDs=%v", columnId, orderedTaskIDs)

	for index, taskID := range orderedTaskIDs {
		_, err := tx.Exec(
			context.Background(),
			`UPDATE tasks
            SET position = $1, column_id = $2, updated_at = NOW()
            WHERE id = $3`,
			index,
			columnId,
			taskID,
		)
		if err != nil {
			log.Printf("ReorderTasks error for task %s: %v", taskID, err)
			return err
		}
	}

	return tx.Commit(context.Background())
}
