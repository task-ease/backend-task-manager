package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/request/query"
	"go-postgres-test/internal/response"
	"go-postgres-test/mixins"
	"strconv"
	"time"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (r *taskRepository) GetALlTasks(data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error) {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.GetALlTasksTx(ctx, tx, data)
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
	if err = r.getPrefixTx(ctx, exec, data.WorkspaceId, data.ProjectId, &prefix); err != nil {
		return nil, err
	}

	for i := range tasks {
		tasks[i].Prefix = prefix + "-" + strconv.Itoa(tasks[i].PrefixNumber)
	}

	return tasks, nil
}

func (r *taskRepository) CreateTask(task request.CreateTask) (response.CreateTask, error) {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return response.CreateTask{}, err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.CreateTaskTx(ctx, tx, task)
}

func (r *taskRepository) CreateTaskTx(ctx context.Context, exec entities.Execer, task request.CreateTask) (response.CreateTask, error) {
	id := uuid.New()

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
		time.Now().UTC(),
		time.Now().UTC(),
	)

	var prefix string
	if err = r.getPrefixTx(ctx, exec, task.WorkspaceId, task.ProjectId, &prefix); err != nil {
		return response.CreateTask{}, err
	}

	var prefixNumber int
	if err = exec.QueryRow(ctx, `
		SELECT prefix_number
		FROM tasks
		WHERE id = $1
	`, id).Scan(&prefixNumber); err != nil {
		return response.CreateTask{}, err
	}

	return response.CreateTask{
		Id:     id,
		Prefix: prefix + "-" + strconv.Itoa(prefixNumber),
	}, err
}

func (r *taskRepository) getPrefixTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID, projectId *uuid.UUID, prefix *string) (err error) {
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

func (r *taskRepository) UpdateTaskTitle(taskId uuid.UUID, value string) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE tasks 
		SET title = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) UpdateTaskDescription(taskId uuid.UUID, value string) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE tasks 
		SET description = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) UpdateTaskColumn(taskId uuid.UUID, columnId uuid.UUID) error {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.UpdateTaskColumnTx(ctx, tx, taskId, columnId)
}

func (r *taskRepository) UpdateTaskColumnTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET column_id = $1
		WHERE id = $2
	`, columnId, taskId)

	var isDoneColumn bool
	err = exec.QueryRow(ctx, `
		SELECT is_done
		FROM task_columns_templates
		WHERE id = $1
	`, columnId).Scan(&isDoneColumn)

	if err != nil {
		return err
	}

	if isDoneColumn {
		_, err = exec.Exec(ctx, `
			UPDATE tasks
			SET is_done = true
			WHERE id = $1
		`, taskId)
		if err != nil {
			return err
		}
	} else {
		var isDoneTask bool
		err = exec.QueryRow(ctx, `
			SELECT is_done
			FROM tasks
			WHERE id = $1
		`, taskId).Scan(&isDoneTask)
		if err != nil {
			return err
		}
		if isDoneTask {
			_, err = exec.Exec(ctx, `
				UPDATE tasks
				SET is_done = false
				WHERE id = $1
			`, taskId)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (r *taskRepository) UpdateTaskPriority(taskId uuid.UUID, priority enums.TaskPriorities) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE tasks 
		SET priority = $1
		WHERE id = $2
	`, priority, taskId)
	return err
}

func (r *taskRepository) UpdateTaskAssigned(taskId uuid.UUID, userId uuid.UUID) error {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.UpdateTaskAssignedTx(ctx, tx, taskId, userId)
}

func (r *taskRepository) UpdateTaskAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		DELETE FROM tasks_assignment
		WHERE task_id = $1
	`, taskId)

	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, `
		INSERT INTO tasks_assignment (task_id, user_id, assigned_at)
		VALUES ($1, $2, $3)
	`, taskId, userId, time.Now().UTC())

	return err
}
