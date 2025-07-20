package handlers

import (
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/types/user"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"
)

type ProjectHandler struct {
	uc usecase.ProjectUseCase
}

func NewProjectHandler(uc usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{uc}
}

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()
	protected := router.Group("/project", middleware.JWTMiddleware(authService))

	protected.GET("/:workspaceId", h.GetAllUserProjects)
	protected.GET("/role/:projectId", h.GetUserRole)
	protected.GET("/members/:projectId")

	protected.POST("", h.CreateProject)
	protected.POST("/user", h.AddUserToProject)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
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
	projectId, err := h.uc.CreateProject(userId, input.WorkspaceId, input.Name, input.Prefix)
	if err != nil {
		if err.Error() == "409" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"projectId": projectId})
}

func (h *ProjectHandler) AddUserToProject(c *gin.Context) {
	var input request.AddUserToProject
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.AddUserToProject(userId, input.ProjectId, input.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": true})
}

func (h *ProjectHandler) GetAllUserProjects(c *gin.Context) {
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

	projects, err := h.uc.GetAllUserProjects(userId, workspaceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *ProjectHandler) GetUserRole(c *gin.Context) {
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

	role, err := h.uc.GetUserRole(userId, projectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == user.ProjectRoleNoAccess {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NO_ACCESS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}
