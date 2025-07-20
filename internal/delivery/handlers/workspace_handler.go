package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/types/user"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"
)

type WorkSpaceHandler struct {
	uc *usecase.WorkSpaceUsecase
}

type createWorkSpaceInput struct {
	Name string `json:"name" binding:"required"`
}

type BindInput struct {
	WorkSpaceId string             `json:"workSpaceId"`
	UserId      string             `json:"userId"`
	Role        user.WorkspaceRole `json:"role"`
}

type RemoveUserInput struct {
	UserId      string `json:"userId"`
	WorkspaceId string `json:"workSpaceId"`
}

func NewWorkSpaceHandler(uc *usecase.WorkSpaceUsecase) *WorkSpaceHandler {
	return &WorkSpaceHandler{uc: uc}
}

func (h *WorkSpaceHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	protected := router.Group("/workspace", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-user-workspaces", h.GetAllUserSpaces)
	protected.GET("/get-all-workspace-members/:id", h.GetAllSpaceMembers)
	protected.GET("/has-user-workspace/:id", h.HasUserWorkspace)
	protected.GET("/search-user/:value", h.SearchWorkspaceMember)

	protected.POST("/change-user-role", h.ChangeUserRole)
	protected.POST("/create-workspace", h.CreateWorkSpace)
	protected.POST("/add-user-to-workspace", h.AddUserToWorkSpace)

	protected.DELETE("/remove-user-from-workspace", h.RemoveUser)
}

func (h *WorkSpaceHandler) CreateWorkSpace(c *gin.Context) {
	var input createWorkSpaceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdStr, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := uuid.Parse(userIdStr.(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var workspace = domain.WorkSpace{
		Name:      input.Name,
		CreatorId: userId,
	}

	status, err := h.uc.CreateWorkSpace(workspace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

func (h *WorkSpaceHandler) GetAllUserSpaces(c *gin.Context) {
	userIdStr, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, err := uuid.Parse(userIdStr.(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workSpacesList, err := h.uc.GetAllUserSpaces(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workSpacesList)
}

func (h *WorkSpaceHandler) AddUserToWorkSpace(c *gin.Context) {
	var input BindInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.uc.AddUserToWorkSpace(input.WorkSpaceId, input.UserId, input.Role)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

//TODO добавить функцию которая будет возвращать не отфильтрованный список, либо просто в workspace_handler передавать либо true либо false

func (h *WorkSpaceHandler) GetAllSpaceMembers(c *gin.Context) {
	workSpaceId := c.Param("id")

	id, err := uuid.Parse(workSpaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	members, err := h.uc.GetAllSpaceMembers(id)

	userIdRaw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userIdStr, ok := userIdRaw.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is not a string"})
		return
	}

	parsedUserId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	membersFiltered := make([]domain.MemberUser, 0)
	for _, m := range members {
		if m.ID != parsedUserId {
			membersFiltered = append(membersFiltered, m)
		}
	}

	c.JSON(http.StatusOK, membersFiltered)
}

func (h *WorkSpaceHandler) RemoveUser(c *gin.Context) {
	var input RemoveUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.uc.RemoveUser(input.WorkspaceId, input.UserId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *WorkSpaceHandler) HasUserWorkspace(c *gin.Context) {
	workSpaceId := c.Param("id")

	anyTypeUserId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userId, ok := anyTypeUserId.(string)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id not found"})
		return
	}

	ok, err := h.uc.HasUserWorkspace(userId, workSpaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ok)
}

func (h *WorkSpaceHandler) ChangeUserRole(c *gin.Context) {
	var input BindInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok, err := h.uc.ChangeUserRole(input.WorkSpaceId, input.UserId, input.Role)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ok)
}

func (h *WorkSpaceHandler) SearchWorkspaceMember(c *gin.Context) {
	workspaceId, err := mixins.QueryToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := mixins.ParseUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value := c.Param("value")
	members, err := h.uc.SearchWorkspaceMember(workspaceId, userId, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}
