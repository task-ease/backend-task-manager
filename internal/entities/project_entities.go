package entities

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	WorkSpaceID uuid.UUID `json:"workSpaceId"`
	CreatorID   uuid.UUID `json:"creatorId"`
	Name        string    `json:"name"`
	IsDone      bool      `json:"isDone"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description *string   `json:"description"`
	Prefix      string    `json:"prefix"`
}
