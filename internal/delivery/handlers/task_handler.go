package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"time"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

type RawTask struct {
	ColumnID    uuid.UUID  `json:"columnId"`
	WorkspaceID uuid.UUID  `json:"workspaceId"`
	Title       string     `json:"title"`
	IsFinished  bool       `json:"isFinished"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	Priority    int        `json:"priority"`
	Status      int        `json:"status"`
}

type ReorderRequest struct {
	ColumnID uuid.UUID   `json:"columnId"`
	TaskIDs  []uuid.UUID `json:"taskIds"`
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/task", middleware.JWTMiddleware(authService))

	protected.POST("/column-tmp", h.CreateColumnTemplate)
}

func (h *TaskHandler) CreateColumnTemplate(c *gin.Context) {

}
