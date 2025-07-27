package request

import "github.com/google/uuid"

type UpdateColumnStatusRequest struct {
	Status   bool      `json:"status"`
	ColumnId uuid.UUID `json:"columnId"`
}

type UpdateColumnNameRequest struct {
	ColumnId uuid.UUID `json:"columnId"`
	Name     string    `json:"name"`
}
