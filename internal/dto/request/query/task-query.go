package query

import "github.com/google/uuid"

type TaskLocationQuery struct {
	WorkspaceId uuid.UUID  `form:"workspaceId"`
	ProjectId   *uuid.UUID `form:"projectId"`
	SprintId    *uuid.UUID `form:"sprintId"`
}
