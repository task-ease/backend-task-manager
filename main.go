package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-postgres-test/infrastructure/db"
	"go-postgres-test/internal/delivery/handlers"
	"go-postgres-test/internal/delivery/ws"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/usecase"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env")
	}

	dbPool := db.ConnectDB()
	defer dbPool.Close()

	userRepo := repository.NewUserRepository(dbPool)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUC)

	taskRepo := repository.NewTaskRepository(dbPool)
	taskUC := usecase.NewTaskUseCase(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskUC)

	messageRepo := repository.NewMessageRepository(dbPool, os.Getenv("IMAGE_STORAGE_API_KEY"))
	messageUC := usecase.NewMessageUsecase(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageUC)

	chatRepo := repository.NewChatRepo(dbPool, messageRepo)
	chatUseCase := usecase.NewChatUsecase(chatRepo)
	chatHandler := handlers.NewChatHandler(*chatUseCase)

	workSpaceRepo := repository.NewWorkSpaceRepository(dbPool, taskRepo)
	workSpaceUC := usecase.NewWorkSpaceUsecase(workSpaceRepo)
	workSpaceHandler := handlers.NewWorkSpaceHandler(workSpaceUC)

	wsHandler := ws.NewWebSocketHandler(messageRepo, userRepo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	wsHandler.RegisterRoutes(r)

	api := r.Group("/api")

	userHandler.RegisterRoutes(api)
	taskHandler.RegisterRoutes(api)
	workSpaceHandler.RegisterRoutes(api)
	chatHandler.RegisterRoutes(api)
	messageHandler.RegisterRoutes(api)

	fmt.Println("Running server")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
