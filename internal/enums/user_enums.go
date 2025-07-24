package enums

type ChatRole string

const (
	ChatUser  ChatRole = "USER"
	ChatAdmin ChatRole = "ADMIN"
)

type ProjectRole string

const (
	ProjectRoleCreator  ProjectRole = "CREATOR"
	ProjectRoleEditor   ProjectRole = "EDITOR"
	ProjectRoleViewer   ProjectRole = "VIEWER"
	ProjectRoleAdmin    ProjectRole = "ADMIN"
	ProjectRoleNoAccess ProjectRole = "NO_ACCESS"
)

type UserRoles string

const (
	User  UserRoles = "USER"
	Admin UserRoles = "ADMIN"
)

type WorkspaceRole string

const (
	WorkspaceCreator    WorkspaceRole = "CREATOR"
	WorkspaceAdmin      WorkspaceRole = "ADMIN"
	WorkspaceMember     WorkspaceRole = "MEMBER"
	WorkspaceNotAllowed WorkspaceRole = "NOT_ALLOWED"
)
