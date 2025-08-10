package repository

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type workSpaceRepo struct {
	conn     *pgxpool.Pool
	taskRepo domain.TaskRepository
}

func NewWorkSpaceRepository(conn *pgxpool.Pool, taskRepo domain.TaskRepository) domain.WorkSpaceRepository {
	return &workSpaceRepo{
		conn:     conn,
		taskRepo: taskRepo,
	}
}

func (r *workSpaceRepo) GetWorkspaceName(ctx context.Context, workspaceId uuid.UUID) (string, error) {
	var name string
	if err := r.conn.QueryRow(ctx, `
		SELECT name 
		FROM workspaces
		WHERE id = $1
	`, workspaceId).Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

func (r *workSpaceRepo) GetIdByColumnTemplateIdTx(ctx context.Context, exec entities.Execer, columnTemplateId uuid.UUID) (uuid.UUID, error) {
	var workspaceId uuid.UUID
	if err := exec.QueryRow(ctx, `
		SELECT workspace_id
		FROM task_columns_templates
		WHERE id = $1
	`, columnTemplateId).Scan(&workspaceId); err != nil {
		return uuid.Nil, err
	}
	return workspaceId, nil
}

func (r *workSpaceRepo) CreateWorkSpaceTx(ctx context.Context, exec entities.Execer, workSpace domain.WorkSpace) error {
	_, err := exec.Exec(ctx,
		`INSERT INTO workspaces (id, creator_id, name) 
		VALUES ($1, $2, $3)
	`, workSpace.ID, workSpace.CreatorId, workSpace.Name)
	return err
}

func (r *workSpaceRepo) AddUser(ctx context.Context, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error {
	return r.AddUserTx(ctx, r.conn, workSpaceId, userId, role)
}

func (r *workSpaceRepo) AddUserTx(ctx context.Context, exec entities.Execer, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error {
	_, err := exec.Exec(ctx,
		`INSERT INTO user_workspaces (user_id, workspace_id, role) 
		VALUES ($1, $2, $3)
	`, userId, workSpaceId, role)
	return err
}

func (r *workSpaceRepo) GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]domain.WorkSpace, error) {
	rows, err := r.conn.Query(ctx, `
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

	if rows.Err() != nil {
		return nil, err
	}

	return workspaces, nil
}

func (r *workSpaceRepo) GetAllMembers(ctx context.Context, workSpaceId uuid.UUID) ([]entities.MemberUser, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT u.id, u.username, u.email, u.icon_url, uw.joined_at, uw.role, uw.position
		FROM user_workspaces uw
		JOIN users u ON uw.user_id = u.id
		WHERE uw.workspace_id = $1`,
		workSpaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.MemberUser
	for rows.Next() {
		var user entities.MemberUser
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.UserIconUrl, &user.JoinedAt, &user.Role, &user.Position); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

func (r *workSpaceRepo) RemoveUser(ctx context.Context, workSpaceId, userId uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM user_workspaces 
        WHERE workspace_id = $1 AND user_id = $2
	`, workSpaceId, userId)
	return err
}

func (r *workSpaceRepo) HasUserWorkspaceTx(ctx context.Context, exec entities.Execer, userId, workspaceId uuid.UUID, exists *bool) error {
	return exec.QueryRow(ctx, `
		SELECT EXISTS (
		SELECT 1 FROM user_workspaces
		WHERE user_id = $1 AND workspace_id = $2
	)`, userId, workspaceId).Scan(&exists)
}

func (r *workSpaceRepo) GetUserRole(ctx context.Context, userId, workspaceId uuid.UUID) (enums.WorkspaceRole, error) {
	return r.GetUserRoleTx(ctx, r.conn, userId, workspaceId)
}

func (r *workSpaceRepo) GetUserRoleTx(ctx context.Context, exec entities.Execer, userId, workspaceId uuid.UUID) (enums.WorkspaceRole, error) {
	var role enums.WorkspaceRole
	if err := exec.QueryRow(ctx, `
		SELECT role 
		FROM user_workspaces 
		WHERE user_id = $1 AND workspace_id = $2
	`, userId, workspaceId).Scan(&role); err != nil {
		return enums.WorkspaceNotAllowed, err
	}
	return role, nil
}

func (r *workSpaceRepo) ChangeUserRole(ctx context.Context, workSpaceId, userId uuid.UUID, role enums.WorkspaceRole) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE user_workspaces
		SET role = $1
		WHERE workspace_id = $2 AND user_id = $3
	`, role, workSpaceId, userId)
	return err
}

func (r *workSpaceRepo) SearchWorkspaceMember(ctx context.Context, workSpaceId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error) {
	value = "%" + value + "%"
	rows, err := r.conn.Query(ctx, `
		SELECT u.id, u.username 
		FROM users u
		JOIN user_workspaces uw ON u.id = uw.user_id
		WHERE uw.workspace_id = $1 
		  AND (u.username ILIKE $2 OR u.email ILIKE $2) 
		ORDER BY u.username
		LIMIT 5
	`, workSpaceId, value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var members []response.FindWorkspaceMemberResponse
	for rows.Next() {
		var member response.FindWorkspaceMemberResponse
		if err = rows.Scan(&member.ID, &member.Name); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return members, nil
}
