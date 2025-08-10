package handlers

import (
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/entities"
	"go-postgres-test/internal/middleware"
	"go-postgres-test/internal/usecase"
	"go-postgres-test/mixins"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	authService := auth.NewJWTService()

	router.GET("/users/is-authorized", h.isAuthorized)

	router.POST("/users/log-in", h.logIn)
	router.POST("/users/create-user", h.createUser)

	protected := router.Group("/users", middleware.JWTMiddleware(authService))

	protected.GET("/search-user-by-email/:value", h.searchUserByEmail)
	protected.GET("/user-id", h.getUserId)
	protected.GET("/workspace-role/:workspaceId", h.getWorkspaceUserRole)
}

func (h *UserHandler) isAuthorized(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	authService := auth.NewJWTService()

	userId, err := authService.VerifyToken(token)

	if err != nil {
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, userId)
}

// TODO почти ничем не отличается от CreateUser, возможно сделать общий, и через GET\POST различать разницу

func (h *UserHandler) logIn(c *gin.Context) {
	var user entities.AuthUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := h.uc.LogIn(c.Request.Context(), user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	authService := auth.NewJWTService()
	token, err := authService.GenerateToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 604800, "/", "localhost", false, true)

	c.Status(http.StatusOK)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var user entities.AuthUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := h.uc.CreateUser(c.Request.Context(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authService := auth.NewJWTService()
	token, err := authService.GenerateToken(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 604800, "/", "localhost", false, true)

	c.Status(http.StatusCreated)
}

func (h *UserHandler) searchUserByEmail(c *gin.Context) {
	users, err := h.uc.SearchUserByEmail(c.Request.Context(), c.Param("value"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) getUserId(c *gin.Context) {
	userIdStr, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userIdStr)
}

func (h *UserHandler) getWorkspaceUserRole(c *gin.Context) {
	userId, err := mixins.ParseUserId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspaceId, err := mixins.QueryToUUID(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRole, err := h.uc.GetWorkspaceUserRole(c.Request.Context(), userId, workspaceId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userRole)
}
