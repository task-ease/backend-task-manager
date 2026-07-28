package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepo struct {
	conn *pgxpool.Pool
}

func NewBaseRepo(conn *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{conn: conn}
}

func (r *BaseRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.conn.Begin(ctx)
}
