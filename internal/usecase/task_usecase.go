package usecase

import (
	"context"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/dto/request"
	"go-postgres-test/internal/dto/request/query"
	"go-postgres-test/internal/dto/response"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/helper"
	"go-postgres-test/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TaskUseCase struct {
	taskRepo domain.TaskRepository
	baseRepo *repository.BaseRepo
}

func NewTaskUseCase(taskRepo domain.TaskRepository, baseRepo *repository.BaseRepo) *TaskUseCase {
	return &TaskUseCase{taskRepo, baseRepo}
}

func (uc *TaskUseCase) GetAllTasks(ctx context.Context, data query.TaskLocationQuery) ([]response.GetAllWorkspaceTasks, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) ([]response.GetAllWorkspaceTasks, error) {
		return uc.taskRepo.GetALlTasksTx(ctx, exec, data)
	})
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, task request.CreateTask) (response.CreateTask, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (response.CreateTask, error) {
		return uc.taskRepo.CreateTaskTx(ctx, exec, task, uuid.New())
	})
}

func (uc *TaskUseCase) UpdateTaskTitle(ctx context.Context, taskId uuid.UUID, value string) error {
	return uc.taskRepo.UpdateTaskTitle(ctx, taskId, value)
}

func (uc *TaskUseCase) UpdateTaskDescription(ctx context.Context, taskId uuid.UUID, value string) error {
	return uc.taskRepo.UpdateTaskDescription(ctx, taskId, value)
}

func (uc *TaskUseCase) UpdateTaskColumn(ctx context.Context, dto request.ChangeTaskPositionAndColumn) (float64, error) {
	return helper.WithTx(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) (float64, error) {
		if err := uc.taskRepo.UpdateTaskColumnIdTx(ctx, exec, dto.TaskId, dto.ToColumnId); err != nil {
			return 0, err
		}

		var isColumnDone bool
		if err := uc.taskRepo.IsColumnDoneTx(ctx, exec, dto.ToColumnId, &isColumnDone); err != nil {
			return 0, err
		}

		if isColumnDone {
			if err := uc.taskRepo.UpdateTaskDoneStatusTx(ctx, exec, dto.TaskId, true); err != nil {
				return 0, err
			}
		} else {
			var isTaskDone bool
			if err := uc.taskRepo.IsTaskDoneTx(ctx, exec, dto.TaskId, &isTaskDone); err != nil {
				return 0, err
			}

			if isTaskDone {
				if err := uc.taskRepo.UpdateTaskDoneStatusTx(ctx, exec, dto.TaskId, false); err != nil {
					return 0, err
				}
			}
		}

		return uc.taskRepo.ChangeTaskPositionAndColumnTx(ctx, exec, dto)
	})
}

func (uc *TaskUseCase) UpdateTaskPriority(ctx context.Context, taskId uuid.UUID, priority enums.TaskPriorities) error {
	return uc.taskRepo.UpdateTaskPriority(ctx, taskId, priority)
}

func (uc *TaskUseCase) UpdateTaskAssigned(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	return helper.WithTxVoid(ctx, uc.baseRepo, func(ctx context.Context, exec pgx.Tx) error {
		var hasAssignment bool
		if err := uc.taskRepo.IsTaskHasAssignedTx(ctx, exec, taskId, &hasAssignment); err != nil {
			return err
		}

		if hasAssignment {
			if err := uc.taskRepo.RemoveTaskAssignedTx(ctx, exec, taskId); err != nil {
				return err
			}
		}

		return uc.taskRepo.UpdateTaskAssignedTx(ctx, exec, taskId, userId)
	})
}

func (uc *TaskUseCase) RemoveTaskAssigned(ctx context.Context, taskId uuid.UUID) error {
	return uc.taskRepo.RemoveTaskAssigned(ctx, taskId)
}

func (uc *TaskUseCase) UpdateTaskType(ctx context.Context, taskId uuid.UUID, value enums.TaskTypes) error {
	return uc.taskRepo.UpdateTaskType(ctx, taskId, value)
}
