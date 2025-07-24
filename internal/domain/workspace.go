package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
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
	AddUserToWorkSpace(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error)
	GetAllSpaceMembers(workSpaceId uuid.UUID) ([]MemberUser, error)
	RemoveUser(workSpaceId string, userId string) (bool, error)
	HasUserWorkspace(userId string, workspaceId string) (enums.WorkspaceRole, error)
	ChangeUserRole(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error)
	SearchWorkspaceMember(workSpaceId, userId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error)
}
