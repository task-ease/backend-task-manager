package handlers

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/internal/usecase"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	//authService := auth.NewJWTService()

	//protected := router.Group("/task", middleware.JWTMiddleware(authService))
}
