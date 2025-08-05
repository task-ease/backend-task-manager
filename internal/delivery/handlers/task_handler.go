package handlers

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/task", middleware.JWTMiddleware(authService))

	protected.GET("/workspace/:workspaceId", h.GetWorkSpaceTasks)
	protected.POST("/", h.CreateTask)

	protected.PATCH("/title/:taskId", h.UpdateTaskTitle)
	protected.PATCH("/description/:taskId", h.UpdateTaskDescription)
	protected.PATCH("/column/:taskId", h.UpdateTaskColumns)
	protected.PATCH("/priority/:taskId", h.UpdateTaskPriority)
	protected.PATCH("/assigned/:taskId", h.UpdateTaskAssigned)
}

func (h *TaskHandler) GetWorkSpaceTasks(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskList, err := h.uc.GetWorkSpaceTasks(workspaceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input request.CreateTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.uc.CreateTask(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

func (h *TaskHandler) UpdateTaskTitle(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value := c.Query("value")

	if err = h.uc.UpdateTaskTitle(taskId, value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) UpdateTaskDescription(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value := c.Query("value")

	if err = h.uc.UpdateTaskDescription(taskId, value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) UpdateTaskColumns(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columnId, err := mixins.QueryToUUID(c, "columnId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.uc.UpdateTaskColumn(taskId, columnId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) UpdateTaskPriority(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priorityStr := c.Query("priority")
	priority := enums.TaskPriorities(priorityStr)

	if err = h.uc.UpdateTaskPriority(taskId, priority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) UpdateTaskAssigned(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.QueryToUUID(c, "userId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.uc.UpdateTaskAssigned(taskId, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
