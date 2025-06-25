package usecase

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/domain"
	"log"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) CreateTask(task *domain.Task) (bool, error) {
	return uc.repo.CreateTask(task)
}

func (uc *TaskUseCase) GetAllColumns(workspaceId uuid.UUID) ([]*domain.TaskColumn, error) {
	return uc.repo.GetAllColumns(workspaceId)
}

func (uc *TaskUseCase) GetAllTasks(workspaceId uuid.UUID) ([]*domain.Task, error) {
	return uc.repo.GetAllTasks(workspaceId)
}

func (uc *TaskUseCase) UpdateTask(task *domain.Task) error {
	return uc.repo.UpdateTask(task)
}

func (uc *TaskUseCase) ReorderTasks(columnId uuid.UUID, orderedTaskIDs []uuid.UUID) error {
	return uc.repo.ReorderTasks(columnId, orderedTaskIDs)
}

var statusMap = map[string]int{
	"todo":        0,
	"in-progress": 1,
	"review":      2,
	"rework":      3,
	"done":        4,
}

var reverseStatusMap = map[int]string{
	0: "todo",
	1: "in-progress",
	2: "review",
	3: "rework",
	4: "done",
}

var priorityMap = map[string]int{
	"low":    0,
	"middle": 1,
	"high":   2,
}

var reversePriorityMap = map[int]string{
	0: "low",
	1: "middle",
	2: "high",
}

func StatusStringToInt(status string) (int, bool) {
	val, ok := statusMap[status]
	if !ok {
		log.Printf("Invalid status received: %s, defaulting to 0 (todo)", status)
		return 0, false
	}
	return val, true
}

func StatusIntToString(i int) string {
	if val, ok := reverseStatusMap[i]; ok {
		return val
	}
	log.Printf("Invalid status int: %d, defaulting to todo", i)
	return "todo"
}

func PriorityStringToInt(priority string) (int, bool) {
	val, ok := priorityMap[priority]
	if !ok {
		log.Printf("Invalid priority received: %s, defaulting to 0 (middle)", priority)
		return 0, false
	}
	return val, true
}

func PriorityIntToString(i int) string {
	if val, ok := reversePriorityMap[i]; ok {
		return val
	}
	log.Printf("Invalid priority int: %d, defaulting to middle", i)
	return "middle"
}
