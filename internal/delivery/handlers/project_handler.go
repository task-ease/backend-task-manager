package handlers

import (
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/enums"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
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

	protected.POST("", h.createProject)
	protected.POST("/user", h.addUserToProject)
}

func (h *ProjectHandler) createProject(c *gin.Context) {
	var input request.CreateProject
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)
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
	var input request.AddUserToProject
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.AddUserToProject(c.Request.Context(), input.ProjectId, input.UserId, enums.ProjectRoleEditor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *ProjectHandler) getAllUserProjects(c *gin.Context) {
	userId, err := mixins.ParseUserId(c)
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
	userId, err := mixins.ParseUserId(c)
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

	if role == enums.ProjectRoleNoAccess {
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
	var input request.ChangeUserProjectRole
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
