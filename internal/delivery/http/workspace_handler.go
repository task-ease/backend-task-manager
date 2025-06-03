package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"net/http"
)

type WorkSpaceHandler struct {
	uc *usecase.WorkSpaceUsecase
}

type createWorkSpaceInput struct {
	Name string `json:"name" binding:"required"`
}

type BindInput struct {
	WorkSpaceId string `json:"workSpaceId"`
	UserId      string `json:"userId"`
	Role        string `json:"role"`
}

func NewWorkSpaceHandler(uc *usecase.WorkSpaceUsecase) *WorkSpaceHandler {
	return &WorkSpaceHandler{uc: uc}
}

func (h *WorkSpaceHandler) RegisterRoutes(router *gin.Engine) {
	authService := auth.NewJWTService()

	protected := router.Group("/workspace", middleware.JWTMiddleware(authService))

	protected.GET("/get-all-user-workspaces", h.GetAllUserSpaces)
	protected.GET("/get-all-workspace-members/:id", h.GetAllSpaceMembers)

	protected.POST("/create-workspace", h.CreateWorkSpace)
	protected.POST("/add-user-to-workspace", h.AddUserToWorkSpace)
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

func (h *WorkSpaceHandler) GetAllSpaceMembers(c *gin.Context) {
	workSpaceId := c.Param("id")

	id, err := uuid.Parse(workSpaceId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	members, err := h.uc.GetAllSpaceMembers(id)

	if err != nil {
		fmt.Printf("error while getting all space members %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)

}
