package repository

import (
	"context"
	"errors"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type projectRepo struct {
	conn *pgxpool.Pool
}

func NewProjectRepository(conn *pgxpool.Pool) domain.ProjectRepository {
	return &projectRepo{conn: conn}
}

func (r *projectRepo) CreateProject(ctx context.Context, creatorId, workSpaceId, projectId uuid.UUID, name, prefix string) error {
	return r.CreateProjectTx(ctx, r.conn, creatorId, workSpaceId, projectId, name, prefix)

}

func (r *projectRepo) CreateProjectTx(ctx context.Context, exec entities.Execer, creatorId, workSpaceId, projectId uuid.UUID, name, prefix string) error {
	var exists bool
	err := exec.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM projects WHERE prefix = $1)`, prefix).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("409")
	}

	now := time.Now().UTC()
	_, err = exec.Exec(ctx, `
		INSERT INTO projects (id, workspace_id, creator_id, name, created_at, updated_at, prefix)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, projectId, workSpaceId, creatorId, name, now, now, prefix)
	if err != nil {
		return err
	}
	return nil
}

func (r *projectRepo) AddUserToProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID, role enums.UserRoles) error {
	return r.AddUserToProjectTx(ctx, r.conn, projectId, userId, role)
}

func (r *projectRepo) AddUserToProjectTx(ctx context.Context, exec entities.Execer, projectId uuid.UUID, userId uuid.UUID, role enums.UserRoles) error {
	_, err := exec.Exec(ctx, `
			INSERT INTO project_members (id, project_id, user_id, role, joined_at)
			VALUES ($1, $2, $3, $4, $5)`, uuid.New(), projectId, userId, role, time.Now().UTC())
	return err
}

func (r *projectRepo) GetAllUserProjects(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT p.id, p.name, p.prefix, p.description FROM projects p
		JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id  = $1 AND p.workspace_id = $2`, userId, workspaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []response.GetAllProjects
	for rows.Next() {
		var project response.GetAllProjects
		if err = rows.Scan(&project.ID, &project.Name, &project.Prefix, &project.Description); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectRepo) GetUserRole(ctx context.Context, userId, projectId uuid.UUID) (enums.UserRoles, error) {
	return r.GetUserRoleTx(ctx, r.conn, userId, projectId)
}

func (r *projectRepo) GetUserRoleTx(ctx context.Context, exec entities.Execer, userId, projectId uuid.UUID) (enums.UserRoles, error) {
	var role enums.UserRoles
	if err := exec.QueryRow(ctx, `
		SELECT 
			CASE 
				WHEN uw.role IN ('ADMIN', 'CREATOR') THEN 'ADMIN'
				ELSE COALESCE(pm.role::text, 'NO_ACCESS')
			END AS effective_role
		FROM projects p
		LEFT JOIN user_workspaces uw ON uw.workspace_id = p.workspace_id AND uw.user_id = $1
		LEFT JOIN project_members pm ON p.id = pm.project_id AND pm.user_id = $1
		WHERE p.id = $2
		LIMIT 1
	`, userId, projectId).Scan(&role); err != nil {
		return enums.NoAccess, err
	}

	return role, nil
}

func (r *projectRepo) GetAllProjectMembers(ctx context.Context, projectId uuid.UUID) ([]response.GetAllProjectUsers, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT DISTINCT ON (u.id)
    		u.id,
    		u.username,
    		u.email,
    		uw.role::text AS workspace_role,
    		pm.role::text AS project_role
		FROM users u
        	JOIN project_members pm ON pm.user_id = u.id
        	JOIN projects p ON pm.project_id = p.id
        	JOIN user_workspaces uw ON uw.user_id = u.id AND uw.workspace_id = p.workspace_id
		WHERE p.id = $1
		ORDER BY u.id`, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectUsers []response.GetAllProjectUsers
	for rows.Next() {
		var projectUser response.GetAllProjectUsers
		if err := rows.Scan(&projectUser.ID, &projectUser.Name, &projectUser.Email, &projectUser.WorkspaceRole, &projectUser.ProjectRole); err != nil {
			return nil, err
		}
		projectUsers = append(projectUsers, projectUser)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return projectUsers, nil
}

func (r *projectRepo) ChangeUserRole(ctx context.Context, userId, projectId uuid.UUID, role enums.UserRoles) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE project_members SET role = $1 WHERE user_id = $2 AND project_id = $3`, role, userId, projectId)
	return err
}

func (r *projectRepo) RemoveUserFromProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM project_members
		WHERE user_id = $1 AND project_id = $2
	`, userId, projectId)
	return err
}

func (r *projectRepo) GetIdByDocumentIdTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	if err := exec.QueryRow(ctx, `
		SELECT project_id
		FROM documents
		WHERE id = $1
	`, documentId).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
