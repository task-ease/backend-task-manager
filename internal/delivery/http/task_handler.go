package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"log"
	"net/http"
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

type SetDoneColumnRequest struct {
	WorkspaceID uuid.UUID `json:"workspaceId"`
	ColumnID    uuid.UUID `json:"columnId"`
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
	protected.POST("/mark-column-done", h.MarkColumnAsDone)
	protected.POST("/update-column", h.UpdateColumn)
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

	task.Priority = &rawTask.Priority
	task.Status = &rawTask.Status

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

type MarkDoneColumnRequest struct {
	ColumnID uuid.UUID `json:"columnId"`
	IsDone   bool      `json:"isDone"`
}

func (h *TaskHandler) MarkColumnAsDone(c *gin.Context) {
	var req MarkDoneColumnRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.uc.MarkColumnAsDone(req.ColumnID, req.IsDone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *TaskHandler) UpdateColumn(c *gin.Context) {
	var input struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("Received update request for column: ID=%s, Name=%s, Color=%s\n", input.ID, input.Name, input.Color)

	if err := h.uc.UpdateColumn(input.ID, input.Name, input.Color); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update column"})
		return
	}

	c.Status(http.StatusOK)
}
