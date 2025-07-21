package response

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/types/user"
)

type GetAllProjects struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	Description *string   `json:"description"`
}

type GetAllProjectUsers struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"username"`
	Email         string             `json:"email"`
	ProjectRole   user.ProjectRole   `json:"projectRole"`
	WorkspaceRole user.WorkspaceRole `json:"workspaceRole"`
}
