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
	dbPool := db.ConnectDB()
	defer dbPool.Close()

	userRepo := repository.NewUserRepository(dbPool)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUC)

	workSpaceRepo := repository.NewWorkSpaceRepository(dbPool)
	workSpaceUC := usecase.NewWorkSpaceUsecase(workSpaceRepo)
	workSpaceHandler := http.NewWorkSpaceHandler(workSpaceUC)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	userHandler.RegisterRoutes(r)
	workSpaceHandler.RegisterRoutes(r)

	fmt.Println("Running server")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
