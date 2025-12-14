package main

import (
	"log"

	"github.com/DjMariarty/messenger/internal/repository"
	service "github.com/DjMariarty/messenger/internal/services"
	handler "github.com/DjMariarty/messenger/internal/transport"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Подключаем базу данных (modernc.org/sqlite, полностью на Go)
	db, err := gorm.Open(sqlite.Open("messenger.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Создаем Gin роутер
	router := gin.Default()

	// Настраиваем маршруты для чатов
	setupChatRoutes(router, db)

	// Запускаем сервер
	log.Println("Server started at http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}

func setupChatRoutes(router *gin.Engine, db *gorm.DB) {
	// 1. Создаем репозиторий
	chatRepo := repository.NewChatRepository(db)

	// 2. Создаем сервис
	chatService := service.NewChatService(db, chatRepo)

	// 3. Создаем хендлер
	chatHandler := handler.NewChatHandler(chatService)

	// 4. Настраиваем API группу
	apiGroup := router.Group("/api")

	// 5. Регистрируем эндпоинт
	apiGroup.GET("/chats", chatHandler.GetChats)
}
