package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/types/user"
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
	if err = r.AddUserToProject(id, creatorId, user.ProjectRoleCreator); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *projectRepo) AddUserToProject(projectId uuid.UUID, userId uuid.UUID, role user.ProjectRole) error {
	_, err := r.conn.Exec(context.Background(), `
			INSERT INTO project_members (id, project_id, user_id, role, joined_at)
			VALUES ($1, $2, $3, $4, $5)`, uuid.New(), projectId, userId, role, time.Now().UTC())
	return err
}
