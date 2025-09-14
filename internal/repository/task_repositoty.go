package repository

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/request/query"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/entities"
	"backend-task-manager/internal/enums"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (r *taskRepository) GetAll(ctx context.Context, data query.TaskLocationWithSprint) ([]response.GetAllWorkspaceTasks, error) {
	return r.GetAllTx(ctx, r.conn, data)
}

func (r *taskRepository) GetAllTx(ctx context.Context, exec entities.Execer, data query.TaskLocationWithSprint) ([]response.GetAllWorkspaceTasks, error) {
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
		ORDER BY t.position
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

func (r *taskRepository) CreateNew(ctx context.Context, task request.CreateTask, id uuid.UUID) (response.CreateTask, error) {
	return r.CreateNewTx(ctx, r.conn, task, id)
}

func (r *taskRepository) CreateNewTx(ctx context.Context, exec entities.Execer, task request.CreateTask, workspaceId uuid.UUID) (response.CreateTask, error) {
	now := time.Now()
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
			created_at,
			updated_at,
			position,
		    prefix_number
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 
		    (
		    	SELECT COUNT(*) * 100 FROM tasks
				WHERE workspace_id = $4
					AND project_id IS NOT DISTINCT FROM $5
				    AND sprint_id IS NOT DISTINCT FROM $6
		    ),
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
		workspaceId,
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
	if err = r.GetPrefixTx(ctx, exec, workspaceId, task.ProjectId, &prefix); err != nil {
		return response.CreateTask{}, err
	}

	var prefixNumber int
	if err = r.GetPrefixNumberTx(ctx, exec, workspaceId, &prefixNumber); err != nil {
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

func (r *taskRepository) UpdateTitle(ctx context.Context, taskId uuid.UUID, value string) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE tasks 
		SET title = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) UpdateDescription(ctx context.Context, taskId uuid.UUID, value string) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE tasks 
		SET description = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) UpdateColumnIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, columnId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET column_id = $1
		WHERE id = $2
	`, columnId, taskId)
	return err
}

func (r *taskRepository) IsDoneTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, done *bool) error {
	err := exec.QueryRow(ctx, `
		SELECT is_done
		FROM tasks
		WHERE id = $1
	`, taskId).Scan(done)
	return err
}

func (r *taskRepository) UpdateDoneStatusTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, to bool) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET is_done = $1
		WHERE id = $2
	`, to, taskId)
	return err
}

func (r *taskRepository) UpdatePriority(ctx context.Context, taskId uuid.UUID, priority enums.TaskPriorities) error {
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

func (r *taskRepository) HasAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, has *bool) error {
	err := exec.QueryRow(ctx, `
		SELECT COUNT(*) > 0
		FROM tasks_assignment
		WHERE task_id = $1
	`, taskId).Scan(has)
	return err
}

func (r *taskRepository) UpdateAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, userId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		INSERT INTO tasks_assignment (task_id, user_id, assigned_at)
		VALUES ($1, $2, $3)
	`, taskId, userId, time.Now().UTC())

	return err
}

func (r *taskRepository) RemoveAssigned(ctx context.Context, taskId uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM tasks_assignment
		WHERE task_id = $1
	`, taskId)
	return err
}

func (r *taskRepository) RemoveAssignedTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		DELETE FROM tasks_assignment
		WHERE task_id = $1
	`, taskId)
	return err
}

func (r *taskRepository) ChangePositionAndColumnTx(
	ctx context.Context,
	exec entities.Execer,
	dto request.ChangeTaskPositionAndColumn,
	taskId uuid.UUID,
) (float64, error) {
	var prevPos float64
	if dto.PrevTaskId != nil {
		if err := exec.QueryRow(ctx, `
			SELECT position FROM tasks WHERE id = $1
		`, dto.PrevTaskId).Scan(&prevPos); err != nil {
			return 0, err
		}
	} else {
		prevPos = 0
	}

	var nextPos float64
	var newPos float64
	if err := exec.QueryRow(ctx, `
		SELECT position 
		FROM tasks 
		WHERE position > $1 AND column_id = $2
		ORDER BY position
		LIMIT 1
	`, prevPos, dto.ToColumnId).Scan(&nextPos); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var columnLength int
			if err = exec.QueryRow(ctx, `
				SELECT COUNT(*)
				FROM tasks 
				WHERE column_id = $1
			`, dto.ToColumnId).Scan(&columnLength); err != nil {
				return 0, err
			}

			newPos = float64(columnLength)*10.0 + 10.0
		} else {
			return 0, err
		}
	} else {
		newPos = (nextPos + prevPos) / 2
	}

	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET position = $1, column_id = $2
		WHERE id = $3
	`, newPos, dto.ToColumnId, taskId)

	if err != nil {
		return 0, err
	}

	return newPos, nil
}

func (r *taskRepository) UpdateType(ctx context.Context, taskId uuid.UUID, value enums.TaskTypes) error {
	return r.UpdateTypeTx(ctx, r.conn, taskId, value)
}

func (r *taskRepository) UpdateTypeTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, value enums.TaskTypes) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET type = $1
		WHERE id = $2
	`, value, taskId)
	return err
}

func (r *taskRepository) Search(ctx context.Context, data query.TaskLocationQuery, value string) ([]response.SearchTasks, error) {
	pattern := "%" + value + "%"

	rows, err := r.conn.Query(ctx, `
    SELECT t.id, t.title
    FROM tasks t
    LEFT JOIN projects p ON t.project_id = p.id
    JOIN workspaces w ON t.workspace_id = w.id
    WHERE t.workspace_id = $1
      	AND (t.project_id = $2 OR ($2 IS NULL AND t.project_id IS NULL))
      	AND (
            (CASE WHEN t.project_id IS NOT NULL THEN p.prefix ELSE w.prefix END)
            || '-' || t.prefix_number || ' ' || t.title
          ) ILIKE $3
    LIMIT 10
	`, data.WorkspaceId, data.ProjectId, pattern)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searchTasks []response.SearchTasks
	for rows.Next() {
		var task response.SearchTasks
		if err = rows.Scan(&task.Id, &task.Title); err != nil {
			return nil, err
		}

		searchTasks = append(searchTasks, task)
	}

	return searchTasks, nil
}

func (r *taskRepository) UpdateParentIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, parentId *uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET parent_id = $1
		WHERE id = $2
	`, parentId, taskId)
	return err
}

func (r *taskRepository) GetTypeTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) (enums.TaskTypes, error) {
	var outputType enums.TaskTypes
	if err := exec.QueryRow(ctx, `
		SELECT type
		FROM tasks 
		WHERE id = $1
	`, taskId).Scan(&outputType); err != nil {
		return "", err
	}

	return outputType, nil
}

func (r *taskRepository) UpdateChildrenIdTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID, childrenId *uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE tasks
		SET parent_id = $1
		WHERE id = $2
	`, taskId, childrenId)
	return err
}

func (r *taskRepository) IfExistsParentTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) (bool, error) {
	var exists bool
	if err := exec.QueryRow(ctx, `
		SELECT parent_id IS NOT NULL
		FROM tasks 
		WHERE id = $1
	`, taskId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *taskRepository) GetLocationTx(ctx context.Context, exec entities.Execer, taskId uuid.UUID) (query.TaskLocationQuery, error) {
	var location query.TaskLocationQuery
	if err := exec.QueryRow(ctx, `
		SELECT (workspace_id, project_id)
		FROM tasks
		WHERE id = $1
	`, taskId).Scan(&location); err != nil {
		return query.TaskLocationQuery{}, err
	}
	return location, nil
}

func (r *taskRepository) GetIdByPrefix(ctx context.Context, prefix string, workspaceId uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	prefixParts := strings.Split(prefix, "-")

	if err := r.conn.QueryRow(ctx, `
		SELECT t.id
		FROM tasks t
		LEFT JOIN projects p ON p.prefix = $1 AND p.workspace_id = $2
		WHERE t.workspace_id = $2 AND t.prefix_number = $3
	`, prefixParts[0], workspaceId, prefixParts[1]).Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
