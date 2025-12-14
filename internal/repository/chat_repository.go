package repository

import (
	"github.com/DjMariarty/messenger/internal/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) GetUserChats(userID uint) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error
	return chats, err
}

func (r *ChatRepository) GetLastMessage(chatID uint) (*models.Message, error) {
	var message models.Message
	err := r.db.Where("chat_id = ?", chatID).Order("created_at DESC").First(&message).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &message, err
}
