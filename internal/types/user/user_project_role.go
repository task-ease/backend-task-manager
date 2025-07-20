package user

type ProjectRole string

const (
	ProjectRoleCreator ProjectRole = "CREATOR"
	ProjectRoleEditor  ProjectRole = "EDITOR"
	ProjectRoleViewer  ProjectRole = "VIEWER"
)
