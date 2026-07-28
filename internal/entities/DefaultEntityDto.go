package entities

import "github.com/google/uuid"

type EntityDto struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}
