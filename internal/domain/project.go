package domain

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	WorkSpaceID uuid.UUID `json:"workSpaceId"`
	CreatorID   uuid.UUID `json:"creatorId"`
	Name        string    `json:"name"`
	IsDone      bool      `json:"isDone"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProjectRepository interface {
	CreateProject(creatorId, workSpaceId uuid.UUID, name, prefix string) (uuid.UUID, error)
}
