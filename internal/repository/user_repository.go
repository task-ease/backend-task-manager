package repository

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/dto"
	"backend-task-manager/internal/entities"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct{ conn *pgxpool.Pool }

func NewUserRepository(conn *pgxpool.Pool) domain.UserRepository {
	return &userRepo{conn: conn}
}

func (r *userRepo) CheckIfExistsByEmail(ctx context.Context, email string, exists *bool) error {
	return r.CheckIfExistsByEmailTx(ctx, r.conn, email, exists)
}

func (r *userRepo) CheckIfExistsByEmailTx(ctx context.Context, exec entities.Execer, email string, exists *bool) error {
	return exec.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`, email).Scan(exists)
}

func (r *userRepo) CreateUser(ctx context.Context, user entities.AuthUser, id uuid.UUID, passwordHash []byte) error {
	return r.CreateUserTx(ctx, r.conn, user, id, passwordHash)
}

func (r *userRepo) CreateUserTx(ctx context.Context, exec entities.Execer, user entities.AuthUser, id uuid.UUID, passwordHash []byte) error {
	_, err := exec.Exec(ctx,
		"INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)",
		id,
		user.Username,
		user.Email,
		passwordHash,
	)
	return err
}

func (r *userRepo) GetIdAndPasswordHash(ctx context.Context, email string) (dto.UserIdAndPasswordHash, error) {
	return r.GetIdAndPasswordHashTx(ctx, r.conn, email)
}

func (r *userRepo) GetIdAndPasswordHashTx(ctx context.Context, exec entities.Execer, email string) (dto.UserIdAndPasswordHash, error) {
	var passwordHash string
	var userId uuid.UUID

	err := exec.QueryRow(ctx, `SELECT password_hash, id FROM users WHERE email = $1`, email).Scan(&passwordHash, &userId)

	if err != nil {
		return dto.UserIdAndPasswordHash{}, err
	}

	return dto.UserIdAndPasswordHash{PasswordHash: passwordHash, ID: userId}, nil
}

func (r *userRepo) SearchUserByEmail(ctx context.Context, value string) ([]entities.User, error) {
	pattern := value + "%"

	rows, err := r.conn.Query(ctx,
		`SELECT email, id, username FROM users WHERE email ILIKE $1 LIMIT 10`,
		pattern)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Email, &user.ID, &user.Username); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) ChangeOnlineStatus(ctx context.Context, userId uuid.UUID, status bool) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE users SET is_online = $1 WHERE id = $2`, status, userId)
	return err
}

func (r *userRepo) GetEmailByUserId(ctx context.Context, userId uuid.UUID) (string, error) {
	var email string
	if err := r.conn.QueryRow(ctx, `SELECT email FROM users WHERE id = $1`, userId).Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}
