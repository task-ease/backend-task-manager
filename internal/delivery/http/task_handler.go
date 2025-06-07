package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"net/http"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

type RawTask struct {
	ColumnID    uuid.UUID      `json:"columnId"`
	WorkspaceID uuid.UUID      `json:"workspaceId"`
	Title       string         `json:"title"`
	IsFinished  bool           `json:"isFinished"`
	Description sql.NullString `json:"description"`
	DueDate     sql.NullTime   `json:"dueDate"`
	Priority    sql.NullInt64  `json:"priority"`
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.Engine) {
	authService := auth.NewJWTService()

	protected := router.Group("/task", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-columns/:id", h.GetAllColumns)
	protected.GET("/get-all-tasks/:id", h.GetAllTasks)

	protected.PATCH("/update-task-title", h.UpdateTaskTitle)

	protected.POST("/create-task", h.CreateTask)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userIdStr, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var task domain.Task
	var rawTask RawTask
	var err error

	task.AuthorID, err = uuid.Parse(userIdStr.(string))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&rawTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ColumnID = rawTask.ColumnID
	task.WorkspaceID = rawTask.WorkspaceID
	task.Title = rawTask.Title
	task.IsFinished = rawTask.IsFinished
	task.Description = rawTask.Description
	task.DueDate = rawTask.DueDate
	task.Priority = rawTask.Priority

	status, err := h.uc.CreateTask(&task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *TaskHandler) GetAllColumns(c *gin.Context) {
	workSpaceId := c.Param("id")

	id, err := uuid.Parse(workSpaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columnsList, err := h.uc.GetAllColumns(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, columnsList)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	workSpaceId := c.Param("id")
	id, err := uuid.Parse(workSpaceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := h.uc.GetAllTasks(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTaskTitle(c *gin.Context) {
	var input struct {
		Title  string    `json:"title"`
		TaskId uuid.UUID `json:"taskId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.UpdateTaskTitle(input.TaskId, input.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
