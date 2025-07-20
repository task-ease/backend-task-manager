package response

import "github.com/google/uuid"

type GetAllProjects struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	Description *string   `json:"description"`
}
