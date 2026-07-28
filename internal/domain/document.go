package domain

import (
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/entities"
	"backend-task-manager/internal/enums"
	"context"

	"github.com/google/uuid"
)

type DocumentRepository interface {
	GetDocument(ctx context.Context, documentId uuid.UUID) (response.GetDocumentResponse, error)
	UpdateDocParent(ctx context.Context, documentId uuid.UUID, parentId *uuid.UUID) error
	UpdateDocumentName(ctx context.Context, documentId uuid.UUID, name string) error
	GetDocumentsByName(ctx context.Context, userId, workspaceId uuid.UUID, name string, limit int) ([]entities.EntityDto, error)
	UpdateDocVisibility(ctx context.Context, documentId uuid.UUID, to enums.DocumentVisibilityTypes) error
	GetIdByNameAndAndWorkspace(ctx context.Context, name string, workspaceId uuid.UUID) (uuid.UUID, error)
	GetAllByUserAndWorkspaceId(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllUserDocumentsResponse, error)
	UpdateDocUserEditPermissions(ctx context.Context, documentId uuid.UUID, req request.UpdateDocUserEditPermissionsRequest) error

	CreateDocTx(ctx context.Context, exec entities.Execer, creatorId, workspaceId uuid.UUID, dto request.CreateDocsRequest) (uuid.UUID, error)
	GetVisibilityTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID) (enums.DocumentVisibilityTypes, error)
	GetUserCanEditTx(ctx context.Context, exec entities.Execer, userId, documentId uuid.UUID) (bool, error)
	UpdateDocContentTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID, content string) error
	CheckIfNameExistsTx(ctx context.Context, exec entities.Execer, name string, workspaceId uuid.UUID) (bool, error)
	UpdateDocVersionsTx(ctx context.Context, exec entities.Execer, documentId, creatorId uuid.UUID, content string) error
}
