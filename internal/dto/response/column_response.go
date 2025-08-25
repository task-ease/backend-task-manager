package response

import "github.com/google/uuid"

type GetWorkspaceColumns struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Color    string    `json:"color"`
	IsDone   bool      `json:"isDone"`
	Position int       `json:"position"`
}
