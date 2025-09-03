package usecase

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/domain/rules"
	"backend-task-manager/internal/dto"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/entities"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/helper"
	"backend-task-manager/internal/repository"
	"backend-task-manager/mixins"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DocumentUsecase struct {
	baseRepo      *repository.BaseRepo
	docRepo       domain.DocumentRepository
	workspaceRepo domain.WorkSpaceRepository
	projectRepo   domain.ProjectRepository
}

func NewDocumentUsecase(
	baseRepo *repository.BaseRepo,
	docRepo domain.DocumentRepository,
	workspaceRepo domain.WorkSpaceRepository,
	projectRepo domain.ProjectRepository,
) *DocumentUsecase {
	return &DocumentUsecase{baseRepo, docRepo, workspaceRepo, projectRepo}
}

func (uc *DocumentUsecase) CheckUserAccess(ctx context.Context, userId, resourceId uuid.UUID) (dto.RolesMiddlewareDto, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (dto.RolesMiddlewareDto, error) {
		var data dto.RolesMiddlewareDto

		visibility, err := uc.docRepo.GetVisibilityTx(ctx, exec, resourceId)

		if err != nil {
			return dto.RolesMiddlewareDto{Role: enums.NoAccess, CanEdit: false}, err
		}

		switch visibility {
		case enums.PrivateDocument:
			data.CanEdit, err = uc.docRepo.GetUserCanEditTx(ctx, exec, userId, resourceId)
		case enums.PublicDocument:
			var role enums.UserRoles
			role, err = uc.workspaceRepo.GetUserRoleByDocumentTx(ctx, exec, userId, resourceId)
			if err != nil {
				return dto.RolesMiddlewareDto{Role: enums.NoAccess, CanEdit: false}, err
			}
			data.CanEdit = mixins.Contains(rules.CanEditWorkspace(), role)
		case enums.ProjectDocument:
			var projectId uuid.UUID
			projectId, err = uc.projectRepo.GetIdByDocumentIdTx(ctx, exec, resourceId)
			if err != nil {
				return dto.RolesMiddlewareDto{Role: enums.NoAccess, CanEdit: false}, err
			}

			var role enums.UserRoles
			role, err = uc.projectRepo.GetUserRoleTx(ctx, exec, userId, projectId)
			if err != nil {
				return dto.RolesMiddlewareDto{Role: enums.NoAccess, CanEdit: false}, err
			}

			data.CanEdit = mixins.Contains(rules.CanEditProject(), role)
		}

		if data.CanEdit {
			data.Role = enums.CanEdit
		} else {
			data.Role = enums.Access
		}

		return data, nil
	})
}

func (uc *DocumentUsecase) CreateDoc(ctx context.Context, creatorId, workspaceId uuid.UUID, dto request.CreateDocsRequest) (uuid.UUID, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (uuid.UUID, error) {
		return uc.CreateDocTx(ctx, exec, creatorId, workspaceId, dto)
	})
}

func (uc *DocumentUsecase) CreateDocTx(ctx context.Context, exec entities.Execer, creatorId, workspaceId uuid.UUID, dto request.CreateDocsRequest) (uuid.UUID, error) {
	exists, err := uc.docRepo.CheckIfNameExistsTx(ctx, exec, dto.Name, workspaceId)

	if err != nil {
		return uuid.Nil, err
	}

	if exists {
		return uuid.Nil, errors.New("name already exists")
	}

	id, err := uc.docRepo.CreateDocTx(ctx, exec, creatorId, workspaceId, dto)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (uc *DocumentUsecase) UpdateDocumentName(ctx context.Context, documentId uuid.UUID, name string) error {
	return uc.docRepo.UpdateDocumentName(ctx, documentId, name)
}

func (uc *DocumentUsecase) UpdateDocParent(ctx context.Context, documentId uuid.UUID, parentId *uuid.UUID) error {
	return uc.docRepo.UpdateDocParent(ctx, documentId, parentId)
}

func (uc *DocumentUsecase) UpdateDocContent(ctx context.Context, documentId, creatorId uuid.UUID, content string) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		if err := uc.docRepo.UpdateDocContentTx(ctx, exec, documentId, content); err != nil {
			return err
		}

		if err := uc.docRepo.UpdateDocVersionsTx(ctx, exec, documentId, creatorId, content); err != nil {
			return err
		}

		return nil
	})
}

func (uc *DocumentUsecase) UpdateDocUserEditPermissions(ctx context.Context, documentId uuid.UUID, req request.UpdateDocUserEditPermissionsRequest) error {
	return uc.docRepo.UpdateDocUserEditPermissions(ctx, documentId, req)
}

func (uc *DocumentUsecase) UpdateDocVisibility(ctx context.Context, documentId uuid.UUID, to enums.DocumentVisibilityTypes) error {
	return uc.docRepo.UpdateDocVisibility(ctx, documentId, to)
}

func (uc *DocumentUsecase) GetAllByUserAndWorkspaceId(ctx context.Context, userId, workspaceId uuid.UUID) ([]response.GetAllUserDocumentsResponse, error) {
	return uc.docRepo.GetAllByUserAndWorkspaceId(ctx, userId, workspaceId)
}

func (uc *DocumentUsecase) GetDocument(ctx context.Context, documentId uuid.UUID) (response.GetDocumentResponse, error) {
	return uc.docRepo.GetDocument(ctx, documentId)
}

func (uc *DocumentUsecase) GetDocumentsByName(ctx context.Context, userId, workspaceId uuid.UUID, name string) ([]entities.EntityDto, error) {
	limit := 5
	return uc.docRepo.GetDocumentsByName(ctx, userId, workspaceId, name, limit)
}
