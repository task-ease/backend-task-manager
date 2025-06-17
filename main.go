package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-postgres-test/infrastructure/db"
	"go-postgres-test/internal/delivery/http"
	"go-postgres-test/internal/delivery/ws"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/usecase"
)

func main() {
	dbPool := db.ConnectDB()
	defer dbPool.Close()

	userRepo := repository.NewUserRepository(dbPool)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUC)

	taskRepo := repository.NewTaskRepository(dbPool)
	taskUC := usecase.NewTaskUseCase(taskRepo)
	taskHandler := http.NewTaskHandler(taskUC)

	messageRepo := repository.NewMessageRepository(dbPool)
	messageUC := usecase.NewMessageUsecase(messageRepo)
	messageHandler := http.NewMessageHandler(messageUC)

	chatRepo := repository.NewChatRepo(dbPool, messageRepo)
	chatUseCase := usecase.NewChatUsecase(chatRepo)
	chatHandler := http.NewChatHandler(*chatUseCase)

	workSpaceRepo := repository.NewWorkSpaceRepository(dbPool, taskRepo)
	workSpaceUC := usecase.NewWorkSpaceUsecase(workSpaceRepo)
	workSpaceHandler := http.NewWorkSpaceHandler(workSpaceUC)

	hub := ws.NewHub()
	wsHandler := http.NewWsHandler(hub)
	go hub.Run()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")

	userHandler.RegisterRoutes(api)
	taskHandler.RegisterRoutes(api)
	workSpaceHandler.RegisterRoutes(api)
	wsHandler.RegisterRoutes(api)
	chatHandler.RegisterRoutes(api)
	messageHandler.RegisterRoutes(api)

	fmt.Println("Running server")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
