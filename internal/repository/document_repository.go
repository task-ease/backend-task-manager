package repository

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/dto/request"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type docRepo struct {
	conn *pgxpool.Pool
}

func NewDocumentRepository(conn *pgxpool.Pool) domain.DocumentRepository {
	return &docRepo{conn: conn}
}

func (r *docRepo) CreateDocTx(ctx context.Context, exec entities.Execer, creatorId, workspaceId uuid.UUID, dto request.CreateDocsRequest) (uuid.UUID, error) {
	id := uuid.New()
	_, err := exec.Exec(ctx, `
		INSERT INTO documents
		(id, name, visibility, workspace_id, project_id, parent_id, creator_id, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, id, dto.Name, dto.Visibility, workspaceId, dto.ProjectId, dto.ParentId, creatorId, time.Now().UTC(), time.Now().UTC())
	return id, err
}

func (r *docRepo) CheckIfNameExistsTx(ctx context.Context, exec entities.Execer, name string, workspaceId uuid.UUID) (bool, error) {
	var exists bool
	if err := exec.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM documents
			WHERE name = $1 AND workspace_id = $2
		)	
	`, name, workspaceId).Scan(&exists); err != nil {
		return exists, err
	}
	return exists, nil
}

func (r *docRepo) UpdateDocContentTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID, content string) error {
	_, err := exec.Exec(ctx, `
		UPDATE documents 
		SET content = $1, updated_at = $2
		WHERE id = $3
	`, content, time.Now().UTC(), documentId)
	return err
}

func (r *docRepo) UpdateDocVersionsTx(ctx context.Context, exec entities.Execer, documentId, creatorId uuid.UUID, content string) error {
	id := uuid.New()
	_, err := exec.Exec(ctx, `
		INSERT INTO document_versions
		(id, document_id, content, created_at, creator_id)
		VALUES ($1, $2, $3, $4)
	`, id, documentId, content, time.Now().UTC(), time.Now().UTC(), creatorId)
	return err
}

func (r *docRepo) GetUserCanEditTx(ctx context.Context, exec entities.Execer, userId, documentId uuid.UUID) (bool, error) {
	var canEdit bool

	if err := exec.QueryRow(ctx, `
		SELECT can_edit
		FROM document_access
		WHERE document_id = $1 AND user_id = $2
	`, documentId, userId).Scan(&canEdit); err != nil {
		return false, err
	}

	return canEdit, nil
}

func (r *docRepo) GetVisibilityTx(ctx context.Context, exec entities.Execer, documentId uuid.UUID) (enums.DocumentVisibilityTypes, error) {
	var visibility enums.DocumentVisibilityTypes
	err := exec.QueryRow(ctx, `
		SELECT visibility
		FROM documents
		WHERE id = $1
	`, documentId).Scan(&visibility)
	return visibility, err
}

func (r *docRepo) UpdateDocumentName(ctx context.Context, documentId uuid.UUID, name string) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE documents
		SET name = $1
		WHERE id = $2
	`, name, documentId)
	return err
}

func (r *docRepo) UpdateDocParent(ctx context.Context, documentId uuid.UUID, parentId *uuid.UUID) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE documents 
		SET parent_id = $1
		WHERE id = $2
	`, parentId, documentId)
	return err
}

func (r *docRepo) UpdateDocUserEditPermissions(ctx context.Context, documentId uuid.UUID, req request.UpdateDocUserEditPermissionsRequest) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE document_access
		SET can_edit = $1
		WHERE document_id = $2 AND user_id = $3
	`, req.To, documentId, req.UserId)
	return err
}

func (r *docRepo) UpdateDocVisibility(ctx context.Context, documentId uuid.UUID, to enums.DocumentVisibilityTypes) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE documents
		SET visibility = $1
		WHERE id = $2
	`, to, documentId)
	return err
}

func (r *docRepo) GetAllByUserAndWorkspaceId(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllUserDocumentsResponse, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT 
		d.id, 
		d.name, 
		d.parent_id,
		COALESCE(
			ARRAY_AGG(d2.id) FILTER (
				WHERE to_tsvector('simple', d2.content) @@ plainto_tsquery('simple', '[[' || d2.name || ']]')
			),
			'{}'
		) as connections_list
		FROM documents d
		LEFT JOIN documents d2 ON TRUE
		LEFT JOIN document_access da ON d.id = da.document_id AND da.user_id = $2
		LEFT JOIN project_members pm ON d.project_id = pm.project_id
		WHERE d.workspace_id = $1
		  AND (
				d.visibility = 'PUBLIC'
			 OR (d.visibility = 'PRIVATE' AND da.user_id = $2)
			 OR (d.visibility = 'PROJECT' AND pm.user_id = $2)
		  )
		GROUP BY d.id, d.name, d.parent_id;
	`, workspaceId, userId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var documents []response.GetAllUserDocumentsResponse
	for rows.Next() {
		var document response.GetAllUserDocumentsResponse
		if err = rows.Scan(&document.Id, &document.Name, &document.ParentId, &document.ConnectionsList); err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	return documents, nil
}

func (r *docRepo) GetDocument(ctx context.Context, documentId uuid.UUID) (response.GetDocumentResponse, error) {
	var document response.GetDocumentResponse
	if err := r.conn.QueryRow(ctx, `
    SELECT 
        d.id, 
		name, 
		content, 
		visibility, 
		parent_id,
		u.username
    FROM documents d
    JOIN users u ON u.id = d.creator_id
    WHERE d.id = $1
`, documentId).Scan(
		&document.Id,
		&document.Name,
		&document.Content,
		&document.Visibility,
		&document.ParentId,
		&document.CreatorName,
	); err != nil {
		return response.GetDocumentResponse{}, err
	}

	return document, nil
}
