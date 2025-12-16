package main

import (
	"log/slog"
	"os"

	"github.com/DjMariarty/messenger/internal/config"
	"github.com/DjMariarty/messenger/internal/middleware"
	"github.com/DjMariarty/messenger/internal/models"
	"github.com/DjMariarty/messenger/internal/repository"
	"github.com/DjMariarty/messenger/internal/services"
	"github.com/DjMariarty/messenger/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.Message{},
	); err != nil {
		log.Error("migrations failed", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("migrations ok")

	userRepo := repository.NewUserRepository(db, log)
	userService := services.NewUserService(db, userRepo, log)
	userHandler := transport.NewUserHandler(userService, log)

	chatRepo := repository.NewChatRepository(db)
	chatService := services.NewChatService(db, chatRepo)
	chatHandler := transport.NewChatHandler(chatService)

	messageRepo := repository.NewMessageRepository(db, log)
	messageService := services.NewMessageService(messageRepo, log)
	messageHandler := transport.NewMessageHandler(messageService, log)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	auth := router.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.GET("/me", middleware.AuthRequired(), userHandler.Me)
	}

	chats := router.Group("/chats")
	chats.Use(middleware.AuthRequired())
	{
		chats.POST("", chatHandler.CreateChat)
		chats.GET("", chatHandler.GetChats)
	}

	messages := router.Group("/messages")
	messages.Use(middleware.AuthRequired())
	{
		messages.POST("/send", messageHandler.CreateMessage)
		messages.GET("/:user_id", messageHandler.GetMessages)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("server starting", slog.String("port", port))
	if err := router.Run(":" + port); err != nil {
		log.Error("http server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
