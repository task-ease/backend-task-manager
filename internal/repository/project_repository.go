package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"time"
)

type projectRepo struct {
	conn *pgxpool.Pool
}

func NewProjectRepository(conn *pgxpool.Pool) domain.ProjectRepository {
	return &projectRepo{conn: conn}
}

func (r *projectRepo) CreateProject(creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error) {
	var id = uuid.New()
	_, err := r.conn.Exec(context.Background(), `
		INSERT INTO projects (id, workspace_id, creator_id, name, created_at, updated_at, prefix)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, id, workSpaceId, creatorId, name, time.Now().UTC(), time.Now().UTC(), prefix)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
