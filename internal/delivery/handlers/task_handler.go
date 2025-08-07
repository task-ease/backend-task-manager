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

	protected.GET("/", h.GetAllTasks)
	protected.POST("/", h.CreateTask)

	protected.PATCH("/title/:taskId", h.UpdateTaskTitle)
	protected.PATCH("/column/:taskId", h.UpdateTaskColumn)
	protected.PATCH("/priority/:taskId", h.UpdateTaskPriority)
	protected.PATCH("/assigned/:taskId", h.UpdateTaskAssigned)
	protected.PATCH("/description/:taskId", h.UpdateTaskDescription)

	protected.DELETE("/assigned/:taskId", h.RemoveTaskAssigned)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
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

	taskList, err := h.uc.GetAllTasks(queryInput)
	if err != nil {
		fmt.Println(err)
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

	data, err := h.uc.CreateTask(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
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

func (h *TaskHandler) UpdateTaskColumn(c *gin.Context) {
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

	priorityStr := c.Query("value")
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

	userId, err := mixins.QueryToUUID(c, "value")
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

func (h *TaskHandler) RemoveTaskAssigned(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.RemoveTaskAssigned(taskId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
