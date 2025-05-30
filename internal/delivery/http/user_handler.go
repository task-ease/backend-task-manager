package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/auth"
	"go-postgres-test/internal/domain"
	"go-postgres-test/internal/usecase"
	"net/http"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/users/create-user", h.CreateUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	authService := auth.NewJWTService()

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Error while getting user ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.uc.CreateUser(user)

	if err != nil {
		fmt.Println("Error while create user ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := authService.GenerateToken(userID)

	if err != nil {
		fmt.Println("Error while generate token ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 86400, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, "success")
}
