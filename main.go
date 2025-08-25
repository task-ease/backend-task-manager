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
	taskRepo := repository.NewTaskRepository(dbPool)
	chatRepo := repository.NewChatRepo(dbPool)
	columnRepo := repository.NewColumnRepo(dbPool)
	projectRepo := repository.NewProjectRepository(dbPool)
	messageRepo := repository.NewMessageRepository(dbPool, os.Getenv("IMAGE_STORAGE_API_KEY"))
	documentRepo := repository.NewDocumentRepository(dbPool)
	workSpaceRepo := repository.NewWorkSpaceRepository(dbPool, taskRepo)

	userUC := usecase.NewUserUsecase(userRepo, baseRepo)
	taskUC := usecase.NewTaskUseCase(taskRepo, baseRepo)
	chatUC := usecase.NewChatUsecase(chatRepo, messageRepo, baseRepo)
	columnUC := usecase.NewColumnUsecase(columnRepo, workSpaceRepo, baseRepo)
	projectUC := usecase.NewProjectUseCase(projectRepo, baseRepo, userRepo)
	messageUC := usecase.NewMessageUsecase(messageRepo, chatRepo, baseRepo)
	documentUC := usecase.NewDocumentUsecase(baseRepo, documentRepo, workSpaceRepo, projectRepo)
	workSpaceUC := usecase.NewWorkSpaceUsecase(workSpaceRepo, columnUC, baseRepo)

	userHandler := handlers.NewUserHandler(userUC)
	taskHandler := handlers.NewTaskHandler(taskUC)
	chatHandler := handlers.NewChatHandler(chatUC)
	columnHandler := handlers.NewColumnHandler(columnUC)
	messageHandler := handlers.NewMessageHandler(messageUC)
	projectHandler := handlers.NewProjectHandler(projectUC)
	documentHandler := handlers.NewDocumentHandler(documentUC, workSpaceUC, projectUC)
	workSpaceHandler := handlers.NewWorkSpaceHandler(workSpaceUC)

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
	chatHandler.RegisterRoutes(api)
	columnHandler.RegisterRoutes(api)
	messageHandler.RegisterRoutes(api)
	projectHandler.RegisterRoutes(api)
	documentHandler.RegisterRoutes(api)
	workSpaceHandler.RegisterRoutes(api)

	port := ":8080"

	if err = r.Run(port); err != nil {
		fmt.Println("Failed to start server:", err)
		panic(err)
	}
}
