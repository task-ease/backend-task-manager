package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
)

type WorkSpaceUsecase struct {
	repo domain.WorkSpaceRepository
}

func NewWorkSpaceUsecase(repo domain.WorkSpaceRepository) *WorkSpaceUsecase {
	return &WorkSpaceUsecase{repo: repo}
}

func (uc *WorkSpaceUsecase) CreateWorkSpace(workspace domain.WorkSpace) (bool, error) {
	return uc.repo.CreateWorkSpace(workspace)
}

func (uc *WorkSpaceUsecase) GetAllUserSpaces(userId uuid.UUID) ([]domain.WorkSpace, error) {
	return uc.repo.GetAllUserSpaces(userId)
}
