package user

type WorkspaceRole string

const (
	WorkspaceCreator    WorkspaceRole = "CREATOR"
	WorkspaceAdmin      WorkspaceRole = "ADMIN"
	WorkspaceMember     WorkspaceRole = "MEMBER"
	WorkspaceNotAllowed WorkspaceRole = "NOT_ALLOWED"
)
