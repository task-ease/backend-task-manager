package response

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
)

type GetAllProjects struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	Description *string   `json:"description"`
}

type GetAllProjectUsers struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"username"`
	Email         string              `json:"email"`
	ProjectRole   enums.ProjectRole   `json:"projectRole"`
	WorkspaceRole enums.WorkspaceRole `json:"workspaceRole"`
}
