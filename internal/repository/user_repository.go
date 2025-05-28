package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-postgres-test/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) domain.UserRepository {
	return &userRepo{conn: conn}
}

func (r *userRepo) CreateUser(user domain.User) error {
	var exists bool

	err := r.conn.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		user.Email,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("email allready exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	_, err = r.conn.Exec(context.Background(),
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
		user.Username,
		user.Email,
		passwordHash,
	)

	return err
}
