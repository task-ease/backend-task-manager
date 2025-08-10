package main

import (
	"fmt"
	"go-postgres-test/infrastructure/db"
	"go-postgres-test/internal/delivery/handlers"
	"go-postgres-test/internal/delivery/ws"
	"go-postgres-test/internal/repository"
	"go-postgres-test/internal/usecase"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env")
	}

	dbPool := db.ConnectDB()
	defer dbPool.Close()

	baseRepo := repository.NewBaseRepo(dbPool)

	userRepo := repository.NewUserRepository(dbPool)
	userUC := usecase.NewUserUsecase(userRepo, baseRepo)
	userHandler := handlers.NewUserHandler(userUC)

	taskRepo := repository.NewTaskRepository(dbPool)
	taskUC := usecase.NewTaskUseCase(taskRepo, baseRepo)
	taskHandler := handlers.NewTaskHandler(taskUC)

	messageRepo := repository.NewMessageRepository(dbPool, os.Getenv("IMAGE_STORAGE_API_KEY"))
	chatRepo := repository.NewChatRepo(dbPool)

	chatUC := usecase.NewChatUsecase(chatRepo, messageRepo, baseRepo)
	chatHandler := handlers.NewChatHandler(chatUC)

	messageUC := usecase.NewMessageUsecase(messageRepo, chatRepo, baseRepo)
	messageHandler := handlers.NewMessageHandler(messageUC)

	workSpaceRepo := repository.NewWorkSpaceRepository(dbPool, taskRepo)
	columnRepo := repository.NewColumnRepo(dbPool)
	columnUC := usecase.NewColumnUsecase(columnRepo, workSpaceRepo, baseRepo)

	workSpaceUC := usecase.NewWorkSpaceUsecase(workSpaceRepo, columnUC, baseRepo)
	workSpaceHandler := handlers.NewWorkSpaceHandler(workSpaceUC)

	columnHandler := handlers.NewColumnHandler(columnUC)

	projectRepo := repository.NewProjectRepository(dbPool)
	projectUC := usecase.NewProjectUseCase(projectRepo, baseRepo)
	projectHandler := handlers.NewProjectHandler(projectUC)

	wsHandler := ws.NewWebSocketHandler(messageRepo, userRepo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PATCH", "PUT"},
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
	projectHandler.RegisterRoutes(api)
	columnHandler.RegisterRoutes(api)

	fmt.Println("Running server")
	if err = r.Run(":8080"); err != nil {
		panic(err)
	}
}
