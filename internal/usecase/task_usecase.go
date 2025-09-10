package usecase

import (
	"backend-task-manager/internal/domain"
	"backend-task-manager/internal/domain/rules"
	"backend-task-manager/internal/dto"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/request/query"
	"backend-task-manager/internal/dto/response"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/helper"
	"backend-task-manager/internal/repository"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskUseCase struct {
	baseRepo    *repository.BaseRepo
	taskRepo    domain.TaskRepository
	projectUC   *ProjectUseCase
	workspaceUC *WorkSpaceUsecase
}

func NewTaskUseCase(
	baseRepo *repository.BaseRepo,
	taskRepo domain.TaskRepository,
	projectUC *ProjectUseCase,
	workspaceUC *WorkSpaceUsecase,
) *TaskUseCase {
	return &TaskUseCase{
		baseRepo,
		taskRepo,
		projectUC,
		workspaceUC,
	}
}

func (uc *TaskUseCase) CheckUserAccess(ctx context.Context, userId, resourceId uuid.UUID) (dto.RolesMiddlewareDto, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (dto.RolesMiddlewareDto, error) {
		location, err := uc.taskRepo.GetLocationTx(ctx, exec, resourceId)
		if err != nil {
			return rules.SettingsNoAccess, err
		}

		if location.ProjectId != nil {
			data, err := uc.projectUC.CheckUserAccess(ctx, userId, *location.ProjectId)
			if err != nil {
				return rules.SettingsNoAccess, err
			}
			return data, nil
		}

		data, err := uc.workspaceUC.CheckUserAccess(ctx, userId, *location.ProjectId)
		if err != nil {
			return rules.SettingsNoAccess, err
		}
		return data, nil
	})
}

func (uc *TaskUseCase) GetAllTasks(ctx context.Context, data query.TaskLocationWithSprintQuery) ([]response.GetAllWorkspaceTasks, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) ([]response.GetAllWorkspaceTasks, error) {
		return uc.taskRepo.GetAllTx(ctx, exec, data)
	})
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, task request.CreateTask) (response.CreateTask, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (response.CreateTask, error) {
		return uc.taskRepo.CreateNewTx(ctx, exec, task, uuid.New())
	})
}

func (uc *TaskUseCase) UpdateTaskTitle(ctx context.Context, taskId uuid.UUID, value string) error {
	return uc.taskRepo.UpdateTitle(ctx, taskId, value)
}

func (uc *TaskUseCase) UpdateTaskDescription(ctx context.Context, taskId uuid.UUID, value string) error {
	return uc.taskRepo.UpdateDescription(ctx, taskId, value)
}

func (uc *TaskUseCase) UpdateTaskColumn(ctx context.Context, dto request.ChangeTaskPositionAndColumn) (float64, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (float64, error) {
		if err := uc.taskRepo.UpdateColumnIdTx(ctx, exec, dto.TaskId, dto.ToColumnId); err != nil {
			return 0, err
		}

		var isColumnDone bool
		if err := uc.taskRepo.IsColumnDoneTx(ctx, exec, dto.ToColumnId, &isColumnDone); err != nil {
			return 0, err
		}

		if isColumnDone {
			if err := uc.taskRepo.UpdateDoneStatusTx(ctx, exec, dto.TaskId, true); err != nil {
				return 0, err
			}
		} else {
			var isTaskDone bool
			if err := uc.taskRepo.IsDoneTx(ctx, exec, dto.TaskId, &isTaskDone); err != nil {
				return 0, err
			}

			if isTaskDone {
				if err := uc.taskRepo.UpdateDoneStatusTx(ctx, exec, dto.TaskId, false); err != nil {
					return 0, err
				}
			}
		}

		return uc.taskRepo.ChangePositionAndColumnTx(ctx, exec, dto)
	})
}

func (uc *TaskUseCase) UpdateTaskPriority(ctx context.Context, taskId uuid.UUID, priority enums.TaskPriorities) error {
	return uc.taskRepo.UpdatePriority(ctx, taskId, priority)
}

func (uc *TaskUseCase) UpdateTaskAssigned(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		var hasAssignment bool
		if err := uc.taskRepo.HasAssignedTx(ctx, exec, taskId, &hasAssignment); err != nil {
			return err
		}

		if hasAssignment {
			if err := uc.taskRepo.RemoveAssignedTx(ctx, exec, taskId); err != nil {
				return err
			}
		}

		return uc.taskRepo.UpdateAssignedTx(ctx, exec, taskId, userId)
	})
}

func (uc *TaskUseCase) RemoveTaskAssigned(ctx context.Context, taskId uuid.UUID) error {
	return uc.taskRepo.RemoveAssigned(ctx, taskId)
}

func (uc *TaskUseCase) UpdateTaskType(ctx context.Context, taskId uuid.UUID, value enums.TaskTypes) error {
	return uc.taskRepo.UpdateType(ctx, taskId, value)
}

func (uc *TaskUseCase) UpdateTaskParent(ctx context.Context, taskId uuid.UUID, parentId *uuid.UUID) (enums.TaskTypes, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (enums.TaskTypes, error) {
		if err := uc.taskRepo.UpdateParentIdTx(ctx, exec, taskId, parentId); err != nil {
			return "", err
		}

		taskType, err := uc.taskRepo.GetTypeTx(ctx, exec, taskId)

		if err != nil {
			return "", err
		}

		if taskType != enums.TaskTypeEpic {
			if parentId == nil {
				if err = uc.taskRepo.UpdateTypeTx(ctx, exec, taskId, enums.TaskTypeTask); err != nil {
					return "", err
				}
				taskType = enums.TaskTypeTask
			} else {
				if err = uc.taskRepo.UpdateTypeTx(ctx, exec, taskId, enums.TaskTypeSubtask); err != nil {
					return "", err
				}
				taskType = enums.TaskTypeSubtask
			}
		}

		return taskType, nil
	})
}

func (uc *TaskUseCase) UpdateTaskChildren(ctx context.Context, taskId uuid.UUID, childrenId *uuid.UUID) (enums.TaskTypes, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (enums.TaskTypes, error) {
		if err := uc.taskRepo.UpdateChildrenIdTx(ctx, exec, taskId, childrenId); err != nil {
			return "", err
		}

		if childrenId != nil {
			if err := uc.taskRepo.UpdateTypeTx(ctx, exec, taskId, enums.TaskTypeEpic); err != nil {
				return enums.TaskTypeEpic, err
			}
		}

		hasParent, err := uc.taskRepo.IfExistsParentTx(ctx, exec, taskId)
		if err != nil {
			return "", err
		}

		newType := enums.TaskTypeTask
		if hasParent {
			newType = enums.TaskTypeSubtask
		}

		if err := uc.taskRepo.UpdateTypeTx(ctx, exec, taskId, newType); err != nil {
			return "", err
		}

		return newType, nil
	})
}
