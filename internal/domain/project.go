package domain

import (
	"context"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	WorkSpaceID uuid.UUID `json:"workSpaceId"`
	CreatorID   uuid.UUID `json:"creatorId"`
	Name        string    `json:"name"`
	IsDone      bool      `json:"isDone"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description *string   `json:"description"`
	Prefix      string    `json:"prefix"`
}

type ProjectRepository interface {
	CreateProject(creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error)
	AddUserToProject(projectId uuid.UUID, userId uuid.UUID, role enums.ProjectRole) error
	ChangeUserRole(userId, projectId uuid.UUID, role enums.ProjectRole) error
	GetAllUserProjects(userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error)
	GetUserRole(userId, projectId uuid.UUID) (enums.ProjectRole, error)
	GetAllProjectMembers(projectId uuid.UUID) ([]response.GetAllProjectUsers, error)

	AddUserToProjectTx(ctx context.Context, exec entities.Execer, projectId uuid.UUID, userId uuid.UUID, role enums.ProjectRole) error
	CreateProjectTx(ctx context.Context, exec entities.Execer, creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error)
}
