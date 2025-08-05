package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"
	"go-postgres-test/mixins"
	"time"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (r *taskRepository) GetWorkSpaceTasks(workspaceId uuid.UUID) ([]response.GetAllWorkspaceTasks, error) {
	rows, err := r.conn.Query(context.Background(), `
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
	WHERE t.workspace_id = $1 AND project_id IS NULL 
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
	    t.position
	ORDER BY t.position
`, workspaceId)

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
			&task.AssignedTo,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// TODO: префикс нужно получать с количества записей в бд тасок, а не из количества элементов в списке на фронте
// Нужно получать число количества тасок где:
// workspaceId == task.workspaceId, projectId == task.projectId, sprintId == task.sprintId
// (если null projectId или sprintId то и в сравнении будет null, то есть равенство выполниться)

func (r *taskRepository) CreateTask(task request.CreateTask) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO tasks
		(
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
		task.PrefixNumber,
	)

	return id, err
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
