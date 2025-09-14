package handlers

import (
	"backend-task-manager/infrastructure/auth"
	"backend-task-manager/internal/domain/rules"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/dto/request/query"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/middleware"
	"backend-task-manager/internal/usecase"
	"backend-task-manager/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskUC      *usecase.TaskUseCase
	workspaceUC *usecase.WorkSpaceUsecase
}

func NewTaskHandler(
	taskUC *usecase.TaskUseCase,
	workspaceUC *usecase.WorkSpaceUsecase,
) *TaskHandler {
	return &TaskHandler{
		taskUC,
		workspaceUC,
	}
}

//TODO сделать отдельный endpoint на создание тасок для проекта

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	protected := router.Group("/task", middleware.JWTMiddleware(auth.NewJWTService()))

	protected.GET("/:"+string(enums.ParamWorkspace), middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.AllWorkspaceRoles), h.getAllTasks)
	protected.GET("/id/:prefix/:"+string(enums.ParamWorkspace), middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.AllWorkspaceRoles), h.getIdByPrefix)

	protected.POST("/:"+string(enums.ParamWorkspace), middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.AllWorkspaceRoles), h.createTask)

	protected.PATCH("/type/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskType)
	protected.PATCH("/title/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskTitle)
	protected.PATCH("/column/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskColumn)
	protected.PATCH("/priority/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskPriority)
	protected.PATCH("/assigned/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskAssigned)
	protected.PATCH("/relations/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskRelations)
	protected.PATCH("/description/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.updateTaskDescription)

	protected.DELETE("/assigned/:"+string(enums.ParamTask), middleware.AccessMiddleware(enums.ParamTask, h.taskUC, rules.HasAccess), h.removeTaskAssigned)
}

func (h *TaskHandler) getAllTasks(c *gin.Context) {
	var queryInput query.TaskLocationWithSprint
	var err error

	queryInput.WorkspaceId, err = mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInput.ProjectId, err = mixins.QueryToUUIDCanBeNull(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queryInput.SprintId, err = mixins.QueryToUUIDCanBeNull(c, "sprintId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskList, err := h.taskUC.GetAllTasks(c.Request.Context(), queryInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

func (h *TaskHandler) createTask(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request.CreateTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskUC.CreateTask(c.Request.Context(), input, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

func (h *TaskHandler) updateTaskTitle(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.taskUC.UpdateTaskTitle(c.Request.Context(), taskId, c.Query("value")); err != nil {
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

	if err = h.taskUC.UpdateTaskDescription(c.Request.Context(), taskId, c.Query("value")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskColumn(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, string(enums.ParamTask))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request.ChangeTaskPositionAndColumn
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPosition, err := h.taskUC.UpdateTaskColumn(c.Request.Context(), input, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newPosition": newPosition})
}

func (h *TaskHandler) updateTaskPriority(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.taskUC.UpdateTaskPriority(c.Request.Context(), taskId, enums.TaskPriorities(c.Query("value"))); err != nil {
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

	if err = h.taskUC.UpdateTaskAssigned(c.Request.Context(), taskId, userId); err != nil {
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

	if err = h.taskUC.RemoveTaskAssigned(c.Request.Context(), taskId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskType(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, string(enums.ParamTask))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.taskUC.UpdateTaskType(c.Request.Context(), taskId, enums.TaskTypes(c.Query("value"))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TaskHandler) updateTaskRelations(c *gin.Context) {
	taskId, err := mixins.ParamToUUID(c, "taskId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relationId, err := mixins.QueryToUUIDCanBeNull(c, "relationId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newType enums.TaskTypes
	switch enums.UpdateTaskRelations(c.Query("method")) {
	case enums.UpdateTaskParent:
		{
			newType, err = h.taskUC.UpdateTaskParent(c.Request.Context(), taskId, relationId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	case enums.UpdateTaskChildren:
		{
			newType, err = h.taskUC.UpdateTaskChildren(c.Request.Context(), taskId, relationId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"newType": newType})
}

func (h *TaskHandler) search(c *gin.Context) {
	var input query.TaskLocationQuery
	var err error

	input.WorkspaceId, err = mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ProjectId, err = mixins.QueryToUUIDCanBeNull(c, string(enums.ParamProject))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskList, err := h.taskUC.Search(c.Request.Context(), input, c.Query("value"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"taskList": taskList})
}

//TODO позже разделить на два endpoint`а (project\workspace) для улучшения безопасности и скорости

func (h *TaskHandler) getIdByPrefix(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefix := c.Param("prefix")

	id, err := h.taskUC.GetIdByPrefix(c.Request.Context(), prefix, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
