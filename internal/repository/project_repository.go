package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
	"time"
)

type projectRepo struct {
	conn *pgxpool.Pool
}

func NewProjectRepository(conn *pgxpool.Pool) domain.ProjectRepository {
	return &projectRepo{conn: conn}
}

func (r *projectRepo) CreateProject(creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error) {
	var exists bool
	err := r.conn.QueryRow(context.Background(), `
		SELECT EXISTS (SELECT 1 FROM projects WHERE prefix = $1)`, prefix).Scan(&exists)
	if err != nil {
		return uuid.Nil, err
	}

	if exists {
		return uuid.Nil, errors.New("409")
	}

	var id = uuid.New()
	_, err = r.conn.Exec(context.Background(), `
		INSERT INTO projects (id, workspace_id, creator_id, name, created_at, updated_at, prefix)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, id, workSpaceId, creatorId, name, time.Now().UTC(), time.Now().UTC(), prefix)
	if err != nil {
		return uuid.Nil, err
	}
	if err = r.AddUserToProject(id, creatorId, enums.ProjectRoleCreator); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *projectRepo) AddUserToProject(projectId uuid.UUID, userId uuid.UUID, role enums.ProjectRole) error {
	_, err := r.conn.Exec(context.Background(), `
			INSERT INTO project_members (id, project_id, user_id, role, joined_at)
			VALUES ($1, $2, $3, $4, $5)`, uuid.New(), projectId, userId, role, time.Now().UTC())
	return err
}

func (r *projectRepo) GetAllUserProjects(userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error) {
	rows, err := r.conn.Query(context.Background(), `
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

	return projects, nil
}

func (r *projectRepo) GetUserRole(userId, projectId uuid.UUID) (enums.ProjectRole, error) {
	var role enums.ProjectRole
	err := r.conn.QueryRow(context.Background(), `
		SELECT 
			CASE 
				WHEN uw.role IN ('ADMIN', 'CREATOR') THEN 'ADMIN'
				ELSE COALESCE(pm.role::text, 'NO_ACCESS')
			END AS effective_role
		FROM projects p
		LEFT JOIN user_workspaces uw ON uw.workspace_id = p.workspace_id AND uw.user_id = $1
		LEFT JOIN project_members pm ON p.id = pm.project_id AND pm.user_id = $1
		WHERE p.id = $2
		LIMIT 1`, userId, projectId).Scan(&role)

	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *projectRepo) GetAllProjectMembers(projectId uuid.UUID) ([]response.GetAllProjectUsers, error) {
	rows, err := r.conn.Query(context.Background(), `
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
	return projectUsers, nil
}

func (r *projectRepo) ChangeUserRole(userId, projectId uuid.UUID, role enums.ProjectRole) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE project_members SET role = $1 WHERE user_id = $2 AND project_id = $3`, role, userId, projectId)
	return err
}
