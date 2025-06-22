package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct{ conn *pgxpool.Pool }

func NewUserRepository(conn *pgxpool.Pool) domain.UserRepository {
	return &userRepo{conn: conn}
}

func (r *userRepo) CreateUser(user domain.AuthUser) (string, error) {
	id := uuid.New().String()
	var exists bool

	err := r.conn.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		user.Email,
	).Scan(&exists)

	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("email already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	_, err = r.conn.Exec(context.Background(),
		"INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)",
		id,
		user.Username,
		user.Email,
		passwordHash,
	)

	return id, err
}

func (r *userRepo) LogIn(user domain.AuthUser) (uuid.UUID, error) {
	var passwordHash string
	var userId uuid.UUID

	err := r.conn.QueryRow(context.Background(),
		`SELECT password_hash, id FROM users WHERE email = $1`,
		user.Email).Scan(&passwordHash, &userId)

	if err != nil {
		return uuid.Nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))

	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (r *userRepo) SearchUserByEmail(value string) ([]domain.User, error) {
	pattern := value + "%"

	rows, err := r.conn.Query(context.Background(),
		`SELECT email, id, username FROM users WHERE email ILIKE $1 LIMIT 10`,
		pattern)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Email, &user.ID, &user.Username); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) ChangeOnlineStatus(userId uuid.UUID, status bool) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE users SET is_online = $1 WHERE id = $2`, status, userId)
	return err
}
