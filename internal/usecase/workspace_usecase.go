package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/response"
	"go-postgres-test/internal/types/user"
)

type WorkSpaceUsecase struct {
	repo domain.WorkSpaceRepository
}

func NewWorkSpaceUsecase(repo domain.WorkSpaceRepository) *WorkSpaceUsecase {
	return &WorkSpaceUsecase{repo: repo}
}

func (uc *WorkSpaceUsecase) CreateWorkSpace(workspace domain.WorkSpace) (uuid.UUID, error) {
	return uc.repo.CreateWorkSpace(workspace)
}

func (uc *WorkSpaceUsecase) GetAllUserSpaces(userId uuid.UUID) ([]domain.WorkSpace, error) {
	return uc.repo.GetAllUserSpaces(userId)
}

func (uc *WorkSpaceUsecase) AddUserToWorkSpace(workSpaceId string, userId string, role user.WorkspaceRole) (bool, error) {
	return uc.repo.AddUserToWorkSpace(workSpaceId, userId, role)
}

func (uc *WorkSpaceUsecase) GetAllSpaceMembers(workSpaceId uuid.UUID) ([]domain.MemberUser, error) {
	return uc.repo.GetAllSpaceMembers(workSpaceId)
}

func (uc *WorkSpaceUsecase) RemoveUser(workSpaceId string, userId string) (bool, error) {
	return uc.repo.RemoveUser(workSpaceId, userId)
}

func (uc *WorkSpaceUsecase) HasUserWorkspace(userId string, workspaceId string) (bool, error) {
	return uc.repo.HasUserWorkspace(userId, workspaceId)
}

func (uc *WorkSpaceUsecase) ChangeUserRole(workSpaceId string, userId string, role user.WorkspaceRole) (bool, error) {
	return uc.repo.ChangeUserRole(workSpaceId, userId, role)
}

func (uc *WorkSpaceUsecase) SearchWorkspaceMember(workSpaceId, userId uuid.UUID, value string) ([]response.FindWorkspaceMemberResponse, error) {
	return uc.repo.SearchWorkspaceMember(workSpaceId, userId, value)
}
