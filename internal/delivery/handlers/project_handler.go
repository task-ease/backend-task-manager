package handlers

import (
	"backend-task-manager/infrastructure/auth"
	"backend-task-manager/internal/dto/request"
	"backend-task-manager/internal/enums"
	"backend-task-manager/internal/middleware"
	"backend-task-manager/internal/usecase"
	"backend-task-manager/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	uc *usecase.ProjectUseCase
}

func NewProjectHandler(uc *usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{uc}
}

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()
	protected := router.Group("/project", middleware.JWTMiddleware(authService))

	protected.GET("/:workspaceId", h.getAllUserProjects)
	protected.GET("/role/:projectId", h.getUserRole)
	protected.GET("/members/:projectId", h.getAllProjectMembers)

	protected.PATCH("/role", h.changeUserRole)

	protected.POST("/", h.createProject)
	protected.POST("/user", h.addUserToProject)

	protected.DELETE("/user", h.removeUserFromProject)
}

func (h *ProjectHandler) createProject(c *gin.Context) {
	var input request.CreateProject
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseContextUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectId, err := h.uc.CreateProject(c.Request.Context(), userId, input.WorkspaceId, input.Name, input.Prefix)
	if err != nil {
		if err.Error() == "409" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"projectId": projectId})
}

func (h *ProjectHandler) addUserToProject(c *gin.Context) {
	var input request.UserProjectManipulation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := h.uc.AddUserToProject(c.Request.Context(), input.ProjectId, input.UserId, enums.ProjectEditor)

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

	projects, err := h.uc.GetAllUserProjects(c.Request.Context(), userId, workspaceId)
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

	role, err := h.uc.GetUserRole(c.Request.Context(), userId, projectId)
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

	members, err := h.uc.GetAllProjectMembers(c.Request.Context(), projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

func (h *ProjectHandler) changeUserRole(c *gin.Context) {
	var input request.ChangeUserUserRoles
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ChangeUserRole(c.Request.Context(), input.UserId, input.ProjectId, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ProjectHandler) removeUserFromProject(c *gin.Context) {
	var input request.UserProjectManipulation

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.RemoveUserFromProject(c.Request.Context(), input.ProjectId, input.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
