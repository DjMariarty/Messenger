package main

import (
	"log"
	"os"

	"github.com/DjMariarty/messenger/internal/config"
	"github.com/DjMariarty/messenger/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {

	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.Message{},
	); err != nil {
		log.Fatal("Не удалось выполнить миграции %v", err)
	}

	router := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер %v", err)
	}
}
