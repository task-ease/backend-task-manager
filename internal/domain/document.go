package domain

import (
	"context"
	"go-postgres-test/internal/dto/request"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"

	"github.com/google/uuid"
)

type DocumentRepository interface {
	GetDocument(ctx context.Context, documentId uuid.UUID) (response.GetDocumentResponse, error)
	UpdateDocParent(ctx context.Context, documentId uuid.UUID, parentId *uuid.UUID) error
	UpdateDocumentName(ctx context.Context, documentId uuid.UUID, name string) error
	UpdateDocVisibility(ctx context.Context, documentId uuid.UUID, to enums.DocumentVisibilityTypes) error
	UpdateDocUserEditPermissions(ctx context.Context, documentId uuid.UUID, req request.UpdateDocUserEditPermissionsRequest) error
	GetAllByUserAndWorkspaceId(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllUserDocumentsResponse, error)

	CreateDocTx(ctx context.Context, exec entities.Execer, creatorId, workspaceId uuid.UUID, dto request.CreateDocsRequest) (uuid.UUID, error)
	GetVisibilityTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID) (enums.DocumentVisibilityTypes, error)
	GetUserCanEditTx(ctx context.Context, exec entities.Execer, userId, documentId uuid.UUID) (bool, error)
	UpdateDocContentTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID, content string) error
	CheckIfNameExistsTx(ctx context.Context, exec entities.Execer, name string, workspaceId uuid.UUID) (bool, error)
	UpdateDocVersionsTx(ctx context.Context, exec entities.Execer, documentId, creatorId uuid.UUID, content string) error
}
