package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"log"
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
	Priority    int            `json:"priority"`
	Status      int            `json:"status"`
}

type ReorderRequest struct {
	ColumnID uuid.UUID   `json:"columnId"`
	TaskIDs  []uuid.UUID `json:"taskIds"`
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{uc: uc} }

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/task", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-columns/:id", h.GetAllColumns)
	protected.GET("/get-all-tasks/:id", h.GetAllTasks)

	protected.POST("/update-task", h.UpdateTask)
	protected.POST("/create-task", h.CreateTask)
	protected.POST("/reorder-tasks", h.ReorderTasks)
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ColumnID = rawTask.ColumnID
	task.WorkspaceID = rawTask.WorkspaceID
	task.Title = rawTask.Title
	task.IsFinished = rawTask.IsFinished

	if rawTask.Description.Valid {
		task.Description = &rawTask.Description.String
	}

	if rawTask.DueDate.Valid {
		task.DueDate = &rawTask.DueDate.Time
	}

	task.Priority = &rawTask.Priority
	task.Status = &rawTask.Status

	status, err := h.uc.CreateTask(&task)
	if err != nil {
		log.Printf("CreateTask error: %v", err)
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
		log.Println("GetAllTasks error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task domain.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ReorderTasks(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ReorderTasks(req.ColumnID, req.TaskIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tasks, err := h.uc.GetAllTasks(uuid.UUID{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "reordered",
		"tasks":  tasks,
	})
}
