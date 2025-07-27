package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/response"
	"time"
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
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		{
			Id:          uuid.Nil,
			WorkspaceId: id,
			Name:        "In Progress",
			Color:       "#6688cc",
			Position:    1,
			IsRequired:  true,
			IsActive:    true,
			IsDone:      false,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		{
			Id:          uuid.Nil,
			WorkspaceId: id,
			Name:        "Done",
			Color:       "#00BFA6",
			Position:    2,
			IsRequired:  true,
			IsActive:    true,
			IsDone:      true,
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

func (uc *WorkSpaceUsecase) GetAllSpaceMembers(workSpaceId uuid.UUID) ([]domain.MemberUser, error) {
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

func (uc *WorkSpaceUsecase) UpdateColumnStatusRequired(columnId uuid.UUID, status bool) error {
	return uc.repo.UpdateColumnStatusRequired(columnId, status)
}
func (uc *WorkSpaceUsecase) UpdateColumnStatusDone(columnId uuid.UUID) error {
	return uc.repo.UpdateColumnStatusDone(columnId)
}

func (uc *WorkSpaceUsecase) UpdateColumnStatusActive(columnId uuid.UUID, status bool) error {
	return uc.repo.UpdateColumnStatusActive(columnId, status)
}

func (uc *WorkSpaceUsecase) UpdateColumnName(columnId uuid.UUID, name string) error {
	return uc.repo.UpdateColumnName(columnId, name)
}
