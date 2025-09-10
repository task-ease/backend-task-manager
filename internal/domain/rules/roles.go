package rules

import (
	"backend-task-manager/internal/dto"
	"backend-task-manager/internal/enums"
)

var AllWorkspaceRoles = []enums.UserRoles{
	enums.WorkspaceCreator,
	enums.WorkspaceAdmin,
	enums.WorkspaceMember,
}

var CanEditWorkspace = []enums.UserRoles{
	enums.WorkspaceCreator,
	enums.WorkspaceAdmin,
	enums.WorkspaceMember,
}

var AllProjectRoles = []enums.UserRoles{
	enums.ProjectCreator,
	enums.ProjectAdmin,
	enums.ProjectEditor,
	enums.ProjectViewer,
}

var CanEditProject = []enums.UserRoles{
	enums.ProjectCreator,
	enums.ProjectAdmin,
	enums.ProjectEditor,
}

var AllDocumentRoles = []enums.UserRoles{
	enums.Access,
	enums.CanEdit,
}

var CanEditDocument = []enums.UserRoles{
	enums.CanEdit,
}

var SettingsNoAccess = dto.RolesMiddlewareDto{
	Role:    enums.NoAccess,
	CanEdit: false,
}
