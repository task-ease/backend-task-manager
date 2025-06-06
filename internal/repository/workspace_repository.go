package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
)

type workSpaceRepo struct {
	conn *pgxpool.Pool
}

func NewWorkSpaceRepository(conn *pgxpool.Pool) domain.WorkSpaceRepository {
	return &workSpaceRepo{conn: conn}
}

func (r *workSpaceRepo) CreateWorkSpace(workSpace domain.WorkSpace) (uuid.UUID, error) {
	workSpace.ID = uuid.New()

	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO workspaces (id, creator_id, name) VALUES ($1, $2, $3)",
		workSpace.ID,
		workSpace.CreatorId,
		workSpace.Name,
	)

	if err != nil {
		return uuid.Nil, err
	}

	_, err = r.AddUserToWorkSpace(workSpace.ID.String(), workSpace.CreatorId.String(), "admin")

	if err != nil {
		return uuid.Nil, err
	}

	return workSpace.ID, nil
}

func (r *workSpaceRepo) AddUserToWorkSpace(workSpaceId string, userId string, role string) (bool, error) {
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

func (r *workSpaceRepo) GetAllSpaceMembers(workSpaceId uuid.UUID) ([]domain.MemberUser, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT u.id, u.username, u.email, uw.joined_at, uw.role
		FROM user_workspaces uw
		JOIN users u ON uw.user_id = u.id
		WHERE uw.workspace_id = $1`,
		workSpaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.MemberUser
	for rows.Next() {
		var user domain.MemberUser
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.JoinedAt, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *workSpaceRepo) RemoveUser(workSpaceId string, userId string) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		`DELETE FROM user_workspaces WHERE workspace_id = $1 AND user_id = $2`,
		workSpaceId,
		userId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *workSpaceRepo) HasUserWorkspace(userId string, workspaceId string) (bool, error) {
	var exists bool

	err := r.conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM user_workspaces
			WHERE user_id = $1 AND workspace_id = $2
			)`,
		userId, workspaceId).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *workSpaceRepo) ChangeUserRole(workSpaceId string, userId string, role string) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE user_workspaces
				SET role = $1
				WHERE workspace_id = $2 AND user_id = $3`, role, workSpaceId, userId)

	if err != nil {
		return false, err
	}

	return true, nil
}
