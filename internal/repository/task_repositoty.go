package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"
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
		 updated_at
		 )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
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

	return id, err
}
