package rules

import "backend-task-manager/internal/enums"

func AllWorkspaceRoles() []enums.UserRoles {
	return []enums.UserRoles{
		enums.WorkspaceCreator,
		enums.WorkspaceAdmin,
		enums.WorkspaceMember,
	}
}

func CanEditWorkspace() []enums.UserRoles {
	return []enums.UserRoles{
		enums.WorkspaceCreator,
		enums.WorkspaceAdmin,
		enums.WorkspaceMember,
	}
}

func AllProjectRoles() []enums.UserRoles {
	return []enums.UserRoles{
		enums.ProjectCreator,
		enums.ProjectAdmin,
		enums.ProjectEditor,
		enums.ProjectViewer,
	}
}

func CanEditProject() []enums.UserRoles {
	return []enums.UserRoles{
		enums.ProjectCreator,
		enums.ProjectAdmin,
		enums.ProjectEditor,
	}
}

func AllDocumentRoles() []enums.UserRoles {
	return []enums.UserRoles{
		enums.Access,
		enums.CanEdit,
	}
}

func CanEditDocument() []enums.UserRoles {
	return []enums.UserRoles{
		enums.CanEdit,
	}
}
