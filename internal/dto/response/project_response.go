package response

import (
	"backend-task-manager/internal/enums"

	"github.com/google/uuid"
)

type GetAllProjects struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	Description *string   `json:"description"`
}

type GetAllProjectUsers struct {
	ID            uuid.UUID       `json:"id"`
	Name          string          `json:"username"`
	Email         string          `json:"email"`
	ProjectRole   enums.UserRoles `json:"projectRole"`
	WorkspaceRole enums.UserRoles `json:"workspaceRole"`
}
