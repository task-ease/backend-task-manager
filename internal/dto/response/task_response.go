package response

import (
	"backend-task-manager/internal/entities"

	"github.com/google/uuid"
)

type GetAllWorkspaceTasks struct {
	entities.Task
	AssignedTo   []AssignedTaskUser `json:"assignedTo"`
	Prefix       string             `json:"prefix"`
	PrefixNumber int                `json:"prefixNumber"`
}

type AssignedTaskUser struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IconUrl  *string   `json:"iconUrl"`
}

type CreateTask struct {
	Id     uuid.UUID `json:"id"`
	Prefix string    `json:"prefix"`
}
