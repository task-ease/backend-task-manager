package entities

import (
	"backend-task-manager/internal/enums"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	Id          uuid.UUID                     `json:"id"`
	Name        string                        `json:"name"`
	CreatorId   uuid.UUID                     `json:"creatorId"`
	CreatedAt   time.Time                     `json:"createdAt"`
	UpdatedAt   time.Time                     `json:"updatedAt"`
	Content     *string                       `json:"content"`
	WorkspaceId uuid.UUID                     `json:"workspaceId"`
	ProjectId   *uuid.UUID                    `json:"projectId"`
	ParentId    *uuid.UUID                    `json:"parentId"`
	Visibility  enums.DocumentVisibilityTypes `json:"visibility"`
}

type DocumentAccess struct {
	Id         uuid.UUID `json:"id"`
	DocumentId uuid.UUID `json:"documentId"`
	UserId     uuid.UUID `json:"userId"`
	CanEdit    bool      `json:"canEdit"`
}

type DocumentVersions struct {
	Id         uuid.UUID `json:"id"`
	DocumentId uuid.UUID `json:"documentId"`
	Content    *string   `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatorId  uuid.UUID `json:"creatorId"`
}
