package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func ConnectDB() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgres://denis:123456@localhost:5432/taskmanager?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	return pool
}
