package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types/user"
	"time"
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
	AddUserToProject(projectId uuid.UUID, userId uuid.UUID, role user.ProjectRole) error
	GetAllUserProjects(userId, workspaceId uuid.UUID) ([]response.GetAllProjects, error)
}
