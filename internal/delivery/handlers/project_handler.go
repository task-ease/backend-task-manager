package handlers

import (
	"backend-task-manager/infrastructure/auth"
	"backend-task-manager/internal/domain/rules"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/middleware"
	"backend-task-manager/internal/usecase"
	"backend-task-manager/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectUC   *usecase.ProjectUseCase
	workspaceUC *usecase.WorkSpaceUsecase
	taskUC      *usecase.TaskUseCase
}

func NewProjectHandler(
	projectUC *usecase.ProjectUseCase,
	workspaceUC *usecase.WorkSpaceUsecase,
	taskUC *usecase.TaskUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		projectUC,
		workspaceUC,
		taskUC,
	}
}

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()
	protected := router.Group("/project", middleware.JWTMiddleware(authService))

	protected.GET("/:"+string(enums.ParamWorkspace), middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.AllWorkspaceRoles), h.getAllUserProjects)
	protected.GET("/role/:"+string(enums.ParamProject), middleware.AccessMiddleware(enums.ParamProject, h.projectUC, rules.AllProjectRoles), h.getUserRole)
	protected.GET("/id/:"+string(enums.ParamWorkspace), middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.AllWorkspaceRoles), h.getProjectIdByPrefix)
	protected.GET("/members/:"+string(enums.ParamProject), middleware.AccessMiddleware(enums.ParamProject, h.projectUC, rules.AllProjectRoles), h.getAllProjectMembers)

	protected.PATCH("/role/:"+string(enums.ParamProject), middleware.AccessMiddleware(enums.ParamProject, h.projectUC, rules.AllProjectRoles), h.changeUserRole)

	protected.PUT("/user/:"+string(enums.ParamProject), middleware.AccessMiddleware(enums.ParamProject, h.projectUC, rules.CanEditProject), h.addUserToProject)

	protected.POST("/:"+string(enums.ParamWorkspace), middleware.AccessMiddleware(enums.ParamWorkspace, h.workspaceUC, rules.CanEditWorkspace), h.createProject)

	protected.DELETE("/user/:"+string(enums.ParamProject), middleware.AccessMiddleware(enums.ParamProject, h.projectUC, rules.CanEditProject), h.removeUserFromProject)
}

func (h *ProjectHandler) createProject(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, string(enums.ParamWorkspace))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request.CreateProject
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.projectUC.CreateProject(c.Request.Context(), userId, workspaceId, input.Name, input.Prefix); err != nil {
		if err.Error() == "409" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ProjectHandler) addUserToProject(c *gin.Context) {
	projectId, err := mixins.ParamToUUID(c, string(enums.ParamProject))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request.UserId
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := h.projectUC.AddUserToProject(c.Request.Context(), projectId, input.UserId, enums.ProjectEditor)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": email})
}

func (h *ProjectHandler) getAllUserProjects(c *gin.Context) {
	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projects, err := h.projectUC.GetAllUserProjects(c.Request.Context(), userId, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *ProjectHandler) getUserRole(c *gin.Context) {
	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectId, err := mixins.ParamToUUID(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.projectUC.GetUserRole(c.Request.Context(), userId, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == enums.NoAccess {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "NO_ACCESS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *ProjectHandler) getAllProjectMembers(c *gin.Context) {
	projectId, err := mixins.ParamToUUID(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	members, err := h.projectUC.GetAllProjectMembers(c.Request.Context(), projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *ProjectHandler) changeUserRole(c *gin.Context) {
	projectId, err := mixins.ParamToUUID(c, string(enums.ParamProject))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request.ChangeUserUserRoles
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.projectUC.ChangeUserRole(c.Request.Context(), input.UserId, projectId, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ProjectHandler) removeUserFromProject(c *gin.Context) {
	projectId, err := mixins.ParamToUUID(c, string(enums.ParamProject))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input request.UserId
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.projectUC.RemoveUserFromProject(c.Request.Context(), projectId, input.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ProjectHandler) getProjectIdByPrefix(c *gin.Context) {
	prefix := c.Query("prefix")

	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.projectUC.GetProjectIdByPrefix(c.Request.Context(), prefix, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
