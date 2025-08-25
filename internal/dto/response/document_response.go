package response

import (
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type GetAllUserDocumentsResponse struct {
	Id              uuid.UUID   `json:"id"`
	Name            string      `json:"name"`
	ParentId        *uuid.UUID  `json:"parentId"`
	ConnectionsList []uuid.UUID `json:"connectionsList"`
}

type GetDocumentResponse struct {
	Id          uuid.UUID                     `json:"id"`
	Content     *string                       `json:"content"`
	CreatorName string                        `json:"creatorName"`
	Name        string                        `json:"name"`
	ParentId    *uuid.UUID                    `json:"parentId"`
	Visibility  enums.DocumentVisibilityTypes `json:"visibility"`
}
