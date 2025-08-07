package usecase

import (
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/response"
	"time"

	"github.com/google/uuid"
)

type WorkSpaceUsecase struct {
	repo domain.WorkSpaceRepository
}

func NewWorkSpaceUsecase(repo domain.WorkSpaceRepository) *WorkSpaceUsecase {
	return &WorkSpaceUsecase{repo: repo}
}

func (uc *WorkSpaceUsecase) CreateWorkSpace(workspace domain.WorkSpace) (uuid.UUID, error) {
	id, err := uc.repo.CreateWorkSpace(workspace)
	if err != nil {
		return uuid.Nil, err
	}

	columnList := []entities.ColumnTemplate{
		{
			Id:          uuid.Nil,
			WorkspaceId: id,
			Name:        "To Do",
			Color:       "#787878",
			Position:    0,
			IsRequired:  true,
			IsActive:    true,
			IsDone:      false,
			GlobalTasks: true,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		{
			Id:          uuid.Nil,
			WorkspaceId: id,
			Name:        "In Progress",
			Color:       "#3d66b8",
			Position:    10,
			IsRequired:  true,
			IsActive:    true,
			IsDone:      false,
			GlobalTasks: true,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		{
			Id:          uuid.Nil,
			WorkspaceId: id,
			Name:        "Done",
			Color:       "#00BFA6",
			Position:    20,
			IsRequired:  true,
			IsActive:    true,
			IsDone:      true,
			GlobalTasks: true,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
	}

	for _, column := range columnList {
		_, err := uc.repo.CreateColumnTemplate(column)
		if err != nil {
			return uuid.Nil, err
		}
	}

	return id, nil
}

func (uc *WorkSpaceUsecase) GetAllUserSpaces(userId uuid.UUID) ([]domain.WorkSpace, error) {
	return uc.repo.GetAllUserSpaces(userId)
}

func (uc *WorkSpaceUsecase) AddUserToWorkSpace(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error) {
	return uc.repo.AddUserToWorkSpace(workSpaceId, userId, role)
}

func (uc *WorkSpaceUsecase) GetAllSpaceMembers(workSpaceId uuid.UUID) ([]entities.MemberUser, error) {
	return uc.repo.GetAllSpaceMembers(workSpaceId)
}

func (uc *WorkSpaceUsecase) RemoveUser(workSpaceId string, userId string) (bool, error) {
	return uc.repo.RemoveUser(workSpaceId, userId)
}

func (uc *WorkSpaceUsecase) HasUserWorkspace(userId string, workspaceId string) (enums.WorkspaceRole, error) {
	return uc.repo.HasUserWorkspace(userId, workspaceId)
}

func (uc *WorkSpaceUsecase) ChangeUserRole(workSpaceId string, userId string, role enums.WorkspaceRole) (bool, error) {
	return uc.repo.ChangeUserRole(workSpaceId, userId, role)
}

func (uc *WorkSpaceUsecase) SearchWorkspaceMember(workSpaceId, userId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error) {
	return uc.repo.SearchWorkspaceMember(workSpaceId, userId, value)
}

func (uc *WorkSpaceUsecase) GetAllColumnTemplates(workSpaceId uuid.UUID) ([]entities.ColumnTemplate, error) {
	return uc.repo.GetAllColumnTemplates(workSpaceId)
}

func (uc *WorkSpaceUsecase) UpdateColumnTemplateStatusRequired(columnId uuid.UUID, status bool) error {
	return uc.repo.UpdateColumnTemplateStatusRequired(columnId, status)
}
func (uc *WorkSpaceUsecase) UpdateColumnTemplateStatusDone(columnId uuid.UUID) error {
	return uc.repo.UpdateColumnTemplateStatusDone(columnId)
}

func (uc *WorkSpaceUsecase) UpdateColumnTemplateStatusActive(columnId uuid.UUID, status bool) error {
	return uc.repo.UpdateColumnTemplateStatusActive(columnId, status)
}

func (uc *WorkSpaceUsecase) UpdateColumnTemplateStatusGlobalTasks(columnId uuid.UUID, status bool) error {
	return uc.repo.UpdateColumnTemplateStatusGlobalTasks(columnId, status)
}

func (uc *WorkSpaceUsecase) UpdateColumnTemplateName(columnId uuid.UUID, name string) error {
	return uc.repo.UpdateColumnTemplateName(columnId, name)
}

func (uc *WorkSpaceUsecase) CreateColumnTemplate(columnTmpReq request.CreateNewColumnTemplate) (uuid.UUID, error) {
	columnTmp := entities.ColumnTemplate{
		WorkspaceId: columnTmpReq.WorkspaceId,
		Name:        columnTmpReq.Name,
		Color:       columnTmpReq.Color,
		Position:    columnTmpReq.Position,
		IsRequired:  columnTmpReq.IsRequired,
		IsDone:      columnTmpReq.IsDone,
		GlobalTasks: columnTmpReq.GlobalTasks,
	}
	return uc.repo.CreateColumnTemplate(columnTmp)
}

func (uc *WorkSpaceUsecase) RenumberColumnTemplatesPositions(workspaceId uuid.UUID) error {
	return uc.repo.RenumberColumnTemplatesPositions(workspaceId)
}

func (uc *WorkSpaceUsecase) GetWorkSpaceColumns(workspaceId uuid.UUID) ([]response.GetWorkSpaceColumns, error) {
	return uc.repo.GetWorkSpaceColumns(workspaceId)
}

func (uc *WorkSpaceUsecase) UpdateColumnTemplateColor(columnId uuid.UUID, color string) error {
	return uc.repo.UpdateColumnTemplateColor(columnId, color)
}
