package user

type WorkspaceRole string

const (
	workspaceCreator WorkspaceRole = "CREATOR"
	workspaceAdmin   WorkspaceRole = "ADMIN"
	workspaceMember  WorkspaceRole = "MEMBER"
)
