package handlers

import (
	"fmt"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/request/query"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/task", middleware.JWTMiddleware(authService))

	protected.GET("/", h.getAllTasks)
	protected.POST("/", h.createTask)

	protected.PATCH("/title/:taskId", h.updateTaskTitle)
	protected.PATCH("/column/:taskId", h.updateTaskColumn)
	protected.PATCH("/priority/:taskId", h.updateTaskPriority)
	protected.PATCH("/assigned/:taskId", h.updateTaskAssigned)
	protected.PATCH("/description/:taskId", h.updateTaskDescription)

	protected.DELETE("/assigned/:taskId", h.removeTaskAssigned)
}

func (h *TaskHandler) getAllTasks(c *gin.Context) {
	var queryInput query.TaskLocationQuery
	var err error

	queryInput.WorkspaceId, err = mixins.QueryToUUID(c, "workspaceId")
	if err != nil {
		fmt.Println("error while workspace id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInput.ProjectId, err = mixins.QueryToUUIDCanBeNull(c, "projectId")
	if err != nil {
		fmt.Println("error while project id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInput.SprintId, err = mixins.QueryToUUIDCanBeNull(c, "sprintId")
	if err != nil {
		fmt.Println("error while sprint id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskList, err := h.uc.GetAllTasks(c.Request.Context(), queryInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

func (h *TaskHandler) createTask(c *gin.Context) {
	var input request.CreateTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.uc.CreateTask(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, data)
}

func (h *TaskHandler) updateTaskTitle(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value := c.Query("value")

	if err = h.uc.UpdateTaskTitle(c.Request.Context(), taskId, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskDescription(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value := c.Query("value")

	if err = h.uc.UpdateTaskDescription(c.Request.Context(), taskId, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskColumn(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columnId, err := mixins.QueryToUUID(c, "value")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.uc.UpdateTaskColumn(c.Request.Context(), taskId, columnId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskPriority(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority := enums.TaskPriorities(c.Query("value"))

	if err = h.uc.UpdateTaskPriority(c.Request.Context(), taskId, priority); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskAssigned(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.QueryToUUID(c, "value")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err = h.uc.UpdateTaskAssigned(c.Request.Context(), taskId, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) removeTaskAssigned(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.RemoveTaskAssigned(c.Request.Context(), taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
