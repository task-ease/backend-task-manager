package repository

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/dto/request"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type columnRepo struct {
	conn *pgxpool.Pool
}

func NewColumnRepo(conn *pgxpool.Pool) domain.ColumnRepository {
	return &columnRepo{conn}
}

func (r *columnRepo) GetColumns(ctx context.Context, workspaceId uuid.UUID, projectId, sprintId *uuid.UUID) ([]response.GetWorkspaceColumns, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT ct.id, name, color, is_done, position
			FROM task_columns_templates ct
			LEFT JOIN using_task_columns utc ON ct.id = utc.template_id 
			WHERE
			(
			  utc.workspace_id = $1
			  AND (utc.project_id = $2 OR ($2 IS NULL AND utc.project_id IS NULL))
			  AND (utc.sprint_id = $3 OR ($3 IS NULL AND utc.sprint_id IS NULL))
			)
			OR
			(ct.workspace_id = $1 AND ct.is_required IS TRUE)
			ORDER BY position
	`, workspaceId, projectId, sprintId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columnList []response.GetWorkspaceColumns
	for rows.Next() {
		var column response.GetWorkspaceColumns
		if err := rows.Scan(
			&column.Id,
			&column.Name,
			&column.Color,
			&column.IsDone,
			&column.Position,
		); err != nil {
			return nil, err
		}

		columnList = append(columnList, column)
	}

	return columnList, nil
}

func (r *columnRepo) CreateColumnTemplateTx(ctx context.Context, exec entities.Execer, columnTmp request.CreateNewColumnTemplate, id uuid.UUID) error {
	now := time.Now().UTC()
	_, err := exec.Exec(ctx, `
		INSERT INTO task_columns_templates (id, workspace_id, name, color, position, is_required, is_done, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, id, columnTmp.WorkspaceId, columnTmp.Name, columnTmp.Color, columnTmp.Position, columnTmp.IsRequired, columnTmp.IsDone, now, now)
	return err
}

func (r *columnRepo) ClearDoneFlagTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
			UPDATE task_columns_templates
			SET is_done = false
			WHERE workspace_id = $1
		`, workspaceId)
	return err
}

func (r *columnRepo) GetAllColumnTemplates(ctx context.Context, workSpaceId uuid.UUID) ([]entities.ColumnTemplate, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT 
			id, 
		  	name, 
		  	color, 
		  	position, 
		  	is_required, 
		  	is_active, 
		  	created_at, 
		  	updated_at, 
		  	is_done, 
		  	EXISTS (
		  	    SELECT 1
		  	    FROM using_task_columns utc
		  	    WHERE utc.template_id = task_columns_templates.id
		  	    	AND utc.project_id IS NULL
		  	    	AND utc.sprint_id IS NULL
		  	) OR is_required = true AS global_tasks
		FROM task_columns_templates
		WHERE workspace_id = $1
		ORDER BY position;
	`, workSpaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columnList []entities.ColumnTemplate
	for rows.Next() {
		var column entities.ColumnTemplate
		if err := rows.Scan(&column.Id, &column.Name, &column.Color, &column.Position, &column.IsRequired, &column.IsActive, &column.CreatedAt, &column.UpdatedAt, &column.IsDone, &column.GlobalTasks); err != nil {
			return nil, err
		}
		columnList = append(columnList, column)
	}

	return columnList, nil
}

func (r *columnRepo) UpdateColumnTemplateName(ctx context.Context, columnTemplateId uuid.UUID, name string) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE task_columns_templates
		SET name = $1
		WHERE id = $2
	`, name, columnTemplateId)
	return err
}

func (r *columnRepo) UpdateColumnTemplateColor(ctx context.Context, columnTemplateId uuid.UUID, color string) error {
	_, err := r.conn.Exec(ctx, `
			UPDATE task_columns_templates
			SET color = $1
			WHERE id = $2	
		`, color, columnTemplateId)
	return err
}

func (r *columnRepo) UpdateColumnTemplateSetDoneTx(ctx context.Context, exec entities.Execer, columnTemplateId uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		UPDATE task_columns_templates
		SET is_done = true
		WHERE id = $1
	`, columnTemplateId)
	return err
}

func (r *columnRepo) AddColumnTx(ctx context.Context, exec entities.Execer, columnTemplateId, workspaceId uuid.UUID, projectId, sprintId *uuid.UUID) (uuid.UUID, error) {
	now := time.Now().UTC()
	id := uuid.New()
	_, err := exec.Exec(ctx, `
	 	INSERT INTO using_task_columns (id, template_id, workspace_id, project_id, sprint_id, created_at, updated_at)
	 	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, columnTemplateId, workspaceId, projectId, sprintId, now, now)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *columnRepo) RemoveColumnTx(ctx context.Context, exec entities.Execer, id uuid.UUID) error {
	_, err := exec.Exec(ctx, `
		DELETE FROM using_task_columns
		WHERE id = $1
	`, id)
	return err
}

func (r *columnRepo) UpdateColumnTemplateStatusRequired(ctx context.Context, columnTemplateId uuid.UUID, status bool) error {
	return r.UpdateColumnTemplateStatusRequiredTx(ctx, r.conn, columnTemplateId, status)
}

func (r *columnRepo) UpdateColumnTemplateStatusRequiredTx(ctx context.Context, exec entities.Execer, columnTemplateId uuid.UUID, status bool) error {
	_, err := exec.Exec(ctx, `
		UPDATE task_columns_templates
		SET is_required = $1
		WHERE id = $2
	`, status, columnTemplateId)
	return err
}

func (r *columnRepo) UpdateColumnTemplateStatusActive(ctx context.Context, columnId uuid.UUID, status bool) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE task_columns_templates
		SET is_active = $1
		WHERE id = $2
	`, status, columnId)
	return err
}

func (r *columnRepo) RenumberColumnTemplatesPositionsTx(ctx context.Context, exec entities.Execer, ids []uuid.UUID) error {
	for i, id := range ids {
		_, err := exec.Exec(ctx, `
            UPDATE task_columns_templates
            SET position = $1
            WHERE id = $2
        `, -(i + 1), id)
		if err != nil {
			return err
		}
	}

	for i, id := range ids {
		_, err := exec.Exec(ctx, `
            UPDATE task_columns_templates
            SET position = $1
            WHERE id = $2
        `, i*10, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *columnRepo) GetColumnTemplatesIdsByWorkspaceOrderByPositionTx(ctx context.Context, exec entities.Execer, workspaceId uuid.UUID) ([]uuid.UUID, error) {
	rows, err := exec.Query(ctx, `
        SELECT id
        FROM task_columns_templates
        WHERE workspace_id = $1
        ORDER BY position
    `, workspaceId)
	if err != nil {
		return []uuid.UUID{}, err
	}
	defer rows.Close()

	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return []uuid.UUID{}, err
	}

	if rows.Err() != nil {
		return []uuid.UUID{}, err
	}

	return ids, nil
}
