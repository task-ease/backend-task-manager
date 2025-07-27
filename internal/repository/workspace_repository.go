package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
	"time"
)

type workSpaceRepo struct {
	conn     *pgxpool.Pool
	taskRepo domain.TaskRepository
}

func NewWorkSpaceRepository(conn *pgxpool.Pool, taskRepo domain.TaskRepository) domain.WorkSpaceRepository {
	return &workSpaceRepo{
		conn:     conn,
		taskRepo: taskRepo,
	}
}

const createColumnTemplateQuery = `
		INSERT INTO task_columns_templates (
		     id, workspace_id, name, color, position, is_required, is_active, is_done, created_at,  updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

func (r *workSpaceRepo) CreateWorkSpace(workSpace domain.WorkSpace) (uuid.UUID, error) {
	workSpace.ID = uuid.New()

	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO workspaces (id, creator_id, name) VALUES ($1, $2, $3)",
		workSpace.ID,
		workSpace.CreatorId,
		workSpace.Name,
	)

	if err != nil {
		return uuid.Nil, err
	}

	//_, err = r.taskRepo.CreateColumn("Planning", workSpace.ID, 0, "38BDF8")
	//_, err = r.taskRepo.CreateColumn("To do", workSpace.ID, 1, "FACC15")
	//_, err = r.taskRepo.CreateColumn("Done", workSpace.ID, 2, "22C55E")

	_, err = r.AddUserToWorkSpace(workSpace.ID.String(), workSpace.CreatorId.String(), enums.WorkspaceCreator)

	if err != nil {
		return uuid.Nil, err
	}

	return workSpace.ID, nil
}

func (r *workSpaceRepo) AddUserToWorkSpace(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO user_workspaces (user_id, workspace_id, role) VALUES ($1, $2, $3)",
		userId,
		workSpaceId,
		role,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *workSpaceRepo) GetAllUserSpaces(userId uuid.UUID) ([]domain.WorkSpace, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT w.id, w.creator_id, w.name, w.created_at
		FROM user_workspaces uw
		JOIN workspaces w ON uw.workspace_id = w.id
		WHERE uw.user_id = $1
		`, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []domain.WorkSpace
	for rows.Next() {
		var ws domain.WorkSpace
		if err := rows.Scan(&ws.ID, &ws.CreatorId, &ws.Name, &ws.CreatedAt); err != nil {
			return nil, err
		}

		workspaces = append(workspaces, ws)
	}

	return workspaces, nil
}

func (r *workSpaceRepo) GetAllSpaceMembers(workSpaceId uuid.UUID) ([]domain.MemberUser, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT u.id, u.username, u.email, u.user_icon_url, uw.joined_at, uw.role, uw.position
		FROM user_workspaces uw
		JOIN users u ON uw.user_id = u.id
		WHERE uw.workspace_id = $1`,
		workSpaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.MemberUser
	for rows.Next() {
		var user domain.MemberUser
		var pos sql.NullString
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.UserIconUrl, &user.JoinedAt, &user.Role, &pos); err != nil {
			return nil, err
		}

		if pos.Valid {
			user.Position = &pos.String
		} else {
			user.Position = nil
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *workSpaceRepo) RemoveUser(workSpaceId string, userId string) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		`DELETE FROM user_workspaces WHERE workspace_id = $1 AND user_id = $2`,
		workSpaceId,
		userId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *workSpaceRepo) HasUserWorkspace(userId string, workspaceId string) (enums.WorkspaceRole, error) {
	var exists bool

	err := r.conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM user_workspaces
			WHERE user_id = $1 AND workspace_id = $2
			)`,
		userId, workspaceId).Scan(&exists)

	if err != nil {
		return enums.WorkspaceNotAllowed, err
	}

	var role enums.WorkspaceRole
	if err = r.conn.QueryRow(context.Background(), `
		SELECT role FROM user_workspaces WHERE user_id = $1`, userId).Scan(&role); err != nil {
		return enums.WorkspaceNotAllowed, err
	}

	return role, nil
}

func (r *workSpaceRepo) ChangeUserRole(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error) {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE user_workspaces
				SET role = $1
				WHERE workspace_id = $2 AND user_id = $3`, role, workSpaceId, userId)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *workSpaceRepo) SearchWorkspaceMember(workSpaceId, userId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error) {
	value = "%" + value + "%"
	rows, err := r.conn.Query(context.Background(), `
		SELECT u.id, u.username 
		FROM users u
		JOIN user_workspaces uw ON u.id = uw.user_id
		WHERE uw.workspace_id = $1 AND (u.username ILIKE $2 OR u.email ILIKE $2) AND u.id != $3
		ORDER BY u.username
		LIMIT 5`,
		workSpaceId, value, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []response.FindWorkspaceMemberResponse
	for rows.Next() {
		var member response.FindWorkspaceMemberResponse
		if err := rows.Scan(&member.ID, &member.Name); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *workSpaceRepo) CreateColumnTemplate(columnTmp entities.ColumnTemplate) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.conn.Exec(context.Background(), createColumnTemplateQuery,
		id,
		columnTmp.WorkspaceId,
		columnTmp.Name,
		columnTmp.Color,
		columnTmp.Position,
		columnTmp.IsRequired,
		columnTmp.IsActive,
		columnTmp.IsDone,
		time.Now().UTC(),
		time.Now().UTC())
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *workSpaceRepo) GetAllColumnTemplates(workSpaceId uuid.UUID) ([]entities.ColumnTemplate, error) {
	rows, err := r.conn.Query(context.Background(), `
		SELECT id, name, color, position, is_required, is_active, created_at, updated_at, is_done FROM task_columns_templates 
		WHERE workspace_id = $1
		ORDER BY position`, workSpaceId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columnList []entities.ColumnTemplate
	for rows.Next() {
		var column entities.ColumnTemplate
		if err := rows.Scan(&column.Id, &column.Name, &column.Color, &column.Position, &column.IsRequired, &column.IsActive, &column.CreatedAt, &column.UpdatedAt, &column.IsDone); err != nil {
			return nil, err
		}
		columnList = append(columnList, column)
	}

	return columnList, nil
}

func (r *workSpaceRepo) UpdateColumnStatusRequired(columnId uuid.UUID, status bool) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE task_columns_templates
		SET is_required = $1
		WHERE id = $2`, status, columnId)
	return err
}

func (r *workSpaceRepo) UpdateColumnStatusDone(columnId uuid.UUID) error {
	var workspaceId uuid.UUID
	err := r.conn.QueryRow(context.Background(), `
		SELECT workspace_id
		FROM task_columns_templates
		WHERE id = $1
	`, columnId).Scan(&workspaceId)
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(context.Background(), `
		UPDATE task_columns_templates
		SET is_done = false
		WHERE workspace_id = $1
	`, workspaceId)
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(context.Background(), `
		UPDATE task_columns_templates
		SET is_done = true, is_required = true
		WHERE id = $1
	`, columnId)
	return err
}

func (r *workSpaceRepo) UpdateColumnStatusActive(columnId uuid.UUID, status bool) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE task_columns_templates
		SET is_active = $1
		WHERE id = $2
	`, status, columnId)
	return err
}

func (r *workSpaceRepo) UpdateColumnName(columnId uuid.UUID, name string) error {
	_, err := r.conn.Exec(context.Background(), `
		UPDATE task_columns_templates
		SET name = $1
		WHERE id = $2
	`, name, columnId)
	return err
}
