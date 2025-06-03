package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go-postgres-test/internal/domain"
)

type workSpaceRepo struct {
	conn *pgx.Conn
}

func NewWorkSpaceRepository(conn *pgx.Conn) domain.WorkSpaceRepository {
	return &workSpaceRepo{conn: conn}
}

func (r *workSpaceRepo) CreateWorkSpace(workSpace domain.WorkSpace) (bool, error) {
	workSpace.ID = uuid.New()

	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO workspaces (id, creator_id, name) VALUES ($1, $2, $3)",
		workSpace.ID,
		workSpace.CreatorId,
		workSpace.Name,
	)

	if err != nil {
		return false, err
	}

	_, err = r.AddUserToWorkSpace(workSpace.ID, workSpace.CreatorId, "admin")

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *workSpaceRepo) AddUserToWorkSpace(workSpaceId uuid.UUID, userId uuid.UUID, role string) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO user_workspaces (user_id, workspace_id, role) VALUES ($1, $2, $3)",
		userId,
		workSpaceId,
		role,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *workSpaceRepo) GetAllUserSpaces(userId uuid.UUID) ([]domain.WorkSpace, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT w.id, w.creator_id, w.name, w.created_at
		FROM user_workspaces uw
		JOIN workspaces w ON uw.workspace_id = w.id
		WHERE uw.user_id = $1
		`, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []domain.WorkSpace
	for rows.Next() {
		var ws domain.WorkSpace
		if err := rows.Scan(&ws.ID, &ws.CreatorId, &ws.Name, &ws.CreatedAt); err != nil {
			return nil, err
		}

		workspaces = append(workspaces, ws)
	}

	return workspaces, nil
}
