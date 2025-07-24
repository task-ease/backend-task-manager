package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"time"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

const createColumnTemplateQuery = `
		INSERT INTO task_columns_templates (
		     id, workspace_id, name, color, position, is_required, is_active, is_done, created_at,  updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

func (r *taskRepository) CreateColumnTemplate(columnTmp entities.ColumnTemplate) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.conn.Exec(context.Background(), createColumnTemplateQuery,
		id,
		columnTmp.WorkspaceId,
		columnTmp.Name,
		columnTmp.Color,
		columnTmp.Position,
		columnTmp.IsRequired,
		columnTmp.IsActive,
		columnTmp.IsDone,
		time.Now().UTC(),
		time.Now().UTC())
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
