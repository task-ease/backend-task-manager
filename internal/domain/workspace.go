package domain

import (
	"github.com/google/uuid"
	"time"
)

type WorkSpace struct {
	ID        uuid.UUID `json:"id"`
	CreatorId uuid.UUID `json:"creator_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkSpaceRepository interface {
	CreateWorkSpace(workspace WorkSpace) (bool, error)
	GetAllUserSpaces(userId uuid.UUID) ([]WorkSpace, error)
}
