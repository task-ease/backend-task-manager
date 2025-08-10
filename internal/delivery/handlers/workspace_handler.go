package handlers

import (
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/request"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkSpaceHandler struct {
	uc *usecase.WorkSpaceUsecase
}

type createWorkSpaceInput struct {
	Name string `json:"name" binding:"required"`
}

func NewWorkSpaceHandler(uc *usecase.WorkSpaceUsecase) *WorkSpaceHandler {
	return &WorkSpaceHandler{uc: uc}
}

func (h *WorkSpaceHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/workspace", middleware.JWTMiddleware(authService))

	protected.GET("/", h.getAllUserWorkSpaces)
	protected.GET("/role/:workspaceId", h.getUserRole)
	protected.GET("/members/:workspaceId", h.getAllMembers)
	protected.GET("/name/:workspaceId", h.getWorkspaceName)
	protected.GET("/search-user/:value", h.searchWorkspaceMember)
	protected.GET("/has-user-workspace/:workspaceId", h.hasUserWorkspace)

	protected.PUT("/role", h.changeUserRole)

	protected.POST("", h.createWorkSpace)
	protected.POST("/user", h.addUserToWorkSpace)

	protected.DELETE("/user", h.removeUser)
}

func (h *WorkSpaceHandler) getAllUserWorkSpaces(c *gin.Context) {
	userId, err := mixins.ParseUserId(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	workSpacesList, err := h.uc.GetAllByUserId(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workSpacesList)
}

func (h *WorkSpaceHandler) createWorkSpace(c *gin.Context) {
	var input createWorkSpaceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	var workspace = domain.WorkSpace{
		Name:      input.Name,
		CreatorId: userId,
	}

	id, err := h.uc.CreateWorkSpace(c.Request.Context(), workspace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, id)
}

func (h *WorkSpaceHandler) addUserToWorkSpace(c *gin.Context) {
	var input request.AddUserToWorkspace
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.AddUser(c.Request.Context(), input.WorkSpaceId, input.UserId, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WorkSpaceHandler) getAllMembers(c *gin.Context) {
	workSpaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	members, err := h.uc.GetAllMembers(c.Request.Context(), workSpaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *WorkSpaceHandler) removeUser(c *gin.Context) {
	var input request.WorkspaceUserManipulation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.RemoveUser(c.Request.Context(), input.WorkspaceId, input.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WorkSpaceHandler) hasUserWorkspace(c *gin.Context) {
	userId, err := mixins.ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	workSpaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.uc.HasUserWorkspace(c.Request.Context(), userId, workSpaceId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *WorkSpaceHandler) changeUserRole(c *gin.Context) {
	var input request.WorkspaceUserManipulationRole
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.ChangeUserRole(c.Request.Context(), input.WorkspaceId, input.UserId, input.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WorkSpaceHandler) searchWorkspaceMember(c *gin.Context) {
	workspaceId, err := mixins.QueryToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value := c.Param("value")

	members, err := h.uc.SearchWorkspaceMember(c.Request.Context(), workspaceId, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *WorkSpaceHandler) getWorkspaceName(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name, err := h.uc.GetWorkspaceName(c.Request.Context(), workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": name})
}

func (h *WorkSpaceHandler) getUserRole(c *gin.Context) {
	workspaceId, err := mixins.ParamToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	role, err := h.uc.GetUserRole(c.Request.Context(), userId, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}
