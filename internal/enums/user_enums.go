package enums

type UserRoles string

const (
	User     UserRoles = "USER"
	Admin    UserRoles = "ADMIN"
	Access   UserRoles = "ACCESS"
	NoAccess UserRoles = "NO_ACCESS"
	CanEdit  UserRoles = "CAN_EDIT"

	ChatUser  UserRoles = "USER"
	ChatAdmin UserRoles = "ADMIN"

	ProjectCreator UserRoles = "CREATOR"
	ProjectEditor  UserRoles = "EDITOR"
	ProjectViewer  UserRoles = "VIEWER"
	ProjectAdmin   UserRoles = "ADMIN"

	WorkspaceCreator UserRoles = "CREATOR"
	WorkspaceAdmin   UserRoles = "ADMIN"
	WorkspaceMember  UserRoles = "MEMBER"
)
