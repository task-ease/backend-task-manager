package handlers

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"net/http"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/task", middleware.JWTMiddleware(authService))

	protected.POST("/column-tmp", h.CreateColumnTemplate)
}

func (h *TaskHandler) CreateColumnTemplate(c *gin.Context) {
	var input entities.ColumnTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.uc.CreateColumnTemplate(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
