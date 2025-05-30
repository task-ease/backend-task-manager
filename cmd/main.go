package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/db"
	"go-postgres-test/internal/delivery/http"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/usecase"
)

func main() {
	db := db.ConnectDB()
	defer db.Close(nil)

	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUC)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	userHandler.RegisterRoutes(r)

	fmt.Println("Running server")
	r.Run(":8080")
}
