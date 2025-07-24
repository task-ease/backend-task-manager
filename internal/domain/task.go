package domain

import (
	"github.com/google/uuid"
	"go-postgres-test/internal/entities"
)

type TaskRepository interface {
	CreateColumnTemplate(columnTmp entities.ColumnTemplate) (uuid.UUID, error)
}
