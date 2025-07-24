package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types/user"
	"time"
)

type WorkSpace struct {
	ID        uuid.UUID `json:"id"`
	CreatorId uuid.UUID `json:"creator_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkSpaceRepository interface {
	CreateWorkSpace(workspace WorkSpace) (uuid.UUID, error)
	GetAllUserSpaces(userId uuid.UUID) ([]WorkSpace, error)
	AddUserToWorkSpace(workSpaceId string, userId string, role user.WorkspaceRole) (bool, error)
	GetAllSpaceMembers(workSpaceId uuid.UUID) ([]MemberUser, error)
	RemoveUser(workSpaceId string, userId string) (bool, error)
	HasUserWorkspace(userId string, workspaceId string) (user.WorkspaceRole, error)
	ChangeUserRole(workSpaceId string, userId string, role user.WorkspaceRole) (bool, error)
	SearchWorkspaceMember(workSpaceId, userId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error)
}
