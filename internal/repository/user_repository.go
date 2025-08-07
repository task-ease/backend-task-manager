package repository

import (
	"context"
	"fmt"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/mixins"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct{ conn *pgxpool.Pool }

func NewUserRepository(conn *pgxpool.Pool) domain.UserRepository {
	return &userRepo{conn: conn}
}

func (r *userRepo) CreateUser(user entities.AuthUser) (id string, err error) {
	ctx := context.Background()

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = mixins.TXReturn(tx, ctx, err)
	}()

	return r.CreateUserTx(ctx, tx, user)
}

func (r *userRepo) CreateUserTx(ctx context.Context, exec entities.Execer, user entities.AuthUser) (string, error) {
	id := uuid.New().String()
	var exists bool

	err := exec.QueryRow(ctx,
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

	_, err = exec.Exec(ctx,
		"INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4)",
		id,
		user.Username,
		user.Email,
		passwordHash,
	)

	return id, err
}

func (r *userRepo) LogIn(user entities.AuthUser) (uuid.UUID, error) {
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

func (r *userRepo) SearchUserByEmail(value string) ([]entities.User, error) {
	pattern := value + "%"

	rows, err := r.conn.Query(context.Background(),
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

	return users, nil
}

func (r *userRepo) ChangeOnlineStatus(userId uuid.UUID, status bool) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE users SET is_online = $1 WHERE id = $2`, status, userId)
	return err
}

func (r *userRepo) GetWorkspaceUserRole(userID uuid.UUID, workspaceId uuid.UUID) (enums.WorkspaceRole, error) {
	var userRole enums.WorkspaceRole
	err := r.conn.QueryRow(context.Background(),
		`SELECT role FROM user_workspaces WHERE user_id = $1 AND workspace_id = $2`,
		userID, workspaceId).Scan(&userRole)
	return userRole, err
}
