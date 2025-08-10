package repository

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/request/query"
	"go-postgres-test/internal/response"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (r *taskRepository) GetALlTasks(ctx context.Context, data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error) {
	return r.GetALlTasksTx(ctx, r.conn, data)
}

func (r *taskRepository) GetALlTasksTx(ctx context.Context, exec entities.Execer, data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error) {
	rows, err := exec.Query(ctx, `
		SELECT 
			t.id, 
			t.column_id, 
			t.parent_id, 
			t.type, 
			t.title, 
			t.description, 
			t.is_done, 
			t.due_date, 
			t.priority, 
			t.position,
			t.prefix_number,
			COALESCE(
				json_agg(
					json_build_object(
						'id', u.id,
						'username', u.username,
						'iconUrl', u.icon_url
					)
				) FILTER (WHERE u.id IS NOT NULL), '[]'
			) AS assigned_to
		FROM tasks t
		LEFT JOIN tasks_assignment ta ON t.id = ta.task_id
		LEFT JOIN users u ON ta.user_id = u.id
		WHERE t.workspace_id = $1 
		  AND (project_id = $2 OR ($2 IS NULL AND project_id IS NULL))
		  AND (sprint_id = $3 OR ($3 IS NULL AND sprint_id IS NULL))
		GROUP BY 
			t.id, 
			t.column_id, 
			t.parent_id, 
			t.type, 
			t.title, 
			t.description, 
			t.is_done, 
			t.due_date, 
			t.priority, 
			t.position,
			t.prefix_number
		ORDER BY t.prefix_number
	`, data.WorkspaceId, data.ProjectId, data.SprintId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []response.GetAllWorkspaceTasks
	for rows.Next() {
		var task response.GetAllWorkspaceTasks
		if err := rows.Scan(
			&task.Id,
			&task.ColumnId,
			&task.ParentId,
			&task.Type,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&task.DueDate,
			&task.Priority,
			&task.Position,
			&task.PrefixNumber,
			&task.AssignedTo,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	var prefix string
	if err = r.GetPrefixTx(ctx, exec, data.WorkspaceId, data.ProjectId, &prefix); err != nil {
		return nil, err
	}

	for i := range tasks {
		tasks[i].Prefix = prefix + "-" + strconv.Itoa(tasks[i].PrefixNumber)
	}

	return tasks, nil
}

func (r *taskRepository) CreateTask(ctx context.Context, task request.CreateTask, id uuid.UUID) (response.CreateTask, error) {
	return r.CreateTaskTx(ctx, r.conn, task, id)
}

func (r *taskRepository) CreateTaskTx(ctx context.Context, exec entities.Execer, task request.CreateTask, id uuid.UUID) (response.CreateTask, error) {
	now := time.Now()

	_, err := exec.Exec(ctx, `
		INSERT INTO tasks (
			id,
			title,
			column_id,
			workspace_id,
			project_id,
			sprint_id,
			author_id,
			parent_id,
			type,
			is_done,
			due_date,
			priority,
			position,
			created_at,
			updated_at,
			prefix_number
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			(
				SELECT COUNT(*) + 1 FROM tasks
				WHERE workspace_id = $4
					AND project_id IS NOT DISTINCT FROM $5
				    AND sprint_id IS NOT DISTINCT FROM $6
			)
		);
	`,
		id,
		task.Title,
		task.ColumnId,
		task.WorkspaceId,
		task.ProjectId,
		task.SprintId,
		task.AuthorId,
		task.ParentId,
		task.Type,
		task.IsDone,
		task.DueDate,
		task.Priority,
		task.Position,
		now,
		now,
	)

	var prefix string
	if err = r.GetPrefixTx(ctx, exec, task.WorkspaceId, task.ProjectId, &prefix); err != nil {
		return response.CreateTask{}, err
	}

	var prefixNumber int
	if err = r.GetPrefixNumberTx(ctx, exec, task.WorkspaceId, &prefixNumber); err != nil {
		return response.CreateTask{}, err
	}

	return response.CreateTask{
		Id:     id,
		Prefix: prefix + "-" + strconv.Itoa(prefixNumber),
	}, err
}

func (r *taskRepository) GetPrefixTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID, projectId *uuid.UUID, prefix *string) (err error) {
	if projectId != nil {
		if err = exec.QueryRow(ctx, `
			SELECT prefix 
			FROM projects
			WHERE id = $1
		`, projectId).Scan(prefix); err != nil {
			return err
		}
	} else {
		if err = exec.QueryRow(ctx, `
			SELECT prefix 
			FROM workspaces
			WHERE id = $1
		`, workspaceId).Scan(prefix); err != nil {
			return err
		}
	}
	return err
}

func (r *taskRepository) GetPrefixNumberTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, prefixNumber *int) error {
	err := exec.QueryRow(ctx, `
		SELECT prefix_number
		FROM tasks
		WHERE id = $1
	`, taskId).Scan(&prefixNumber)
	return err
}

func (r *taskRepository) UpdateTaskTitle(ctx context.Context, taskId uuid.UUID, value string) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE tasks 
		SET title = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) UpdateTaskDescription(ctx context.Context, taskId uuid.UUID, value string) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE tasks 
		SET description = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) UpdateTaskColumnIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET column_id = $1
		WHERE id = $2
	`, columnId, taskId)
	return err
}

func (r *taskRepository) IsTaskDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error {
	err := exec.QueryRow(ctx, `
		SELECT is_done
		FROM tasks
		WHERE id = $1
	`, taskId).Scan(done)
	return err
}

func (r *taskRepository) UpdateTaskDoneStatusTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, to bool) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET is_done = $1
		WHERE id = $2
	`, to, taskId)
	return err
}

func (r *taskRepository) UpdateTaskPriority(ctx context.Context, taskId uuid.UUID, priority enums.TaskPriorities) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE tasks 
		SET priority = $1
		WHERE id = $2
	`, priority, taskId)
	return err
}

func (r *taskRepository) IsColumnDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error {
	err := exec.QueryRow(ctx, `
		SELECT is_done 
		FROM task_columns_templates
		WHERE id = $1
	`, taskId).Scan(done)
	return err
}

func (r *taskRepository) IsTaskHasAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, has *bool) error {
	err := exec.QueryRow(ctx, `
		SELECT COUNT(*) > 0
		FROM tasks_assignment
		WHERE task_id = $1
	`, taskId).Scan(has)
	return err
}

func (r *taskRepository) UpdateTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO tasks_assignment (task_id, user_id, assigned_at)
		VALUES ($1, $2, $3)
	`, taskId, userId, time.Now().UTC())

	return err
}

func (r *taskRepository) RemoveTaskAssigned(ctx context.Context, taskId uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM tasks_assignment
		WHERE task_id = $1
	`, taskId)
	return err
}

func (r *taskRepository) RemoveTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		DELETE FROM tasks_assignment
		WHERE task_id = $1
	`, taskId)
	return err
}
