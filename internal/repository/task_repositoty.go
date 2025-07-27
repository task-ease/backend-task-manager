package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
)

type taskRepository struct{ conn *pgxpool.Pool }

func NewTaskRepository(conn *pgxpool.Pool) domain.TaskRepository {
	return &taskRepository{conn: conn}
}
