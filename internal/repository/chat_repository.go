package repository

import (
	"errors"

	"github.com/DjMariarty/messenger/internal/models"
	"gorm.io/gorm"
)

type ChatRepository interface {
	FindByUsers(user1ID, user2ID uint) (*models.Chat, error)
	Create(chat *models.Chat) error
	GetUserChats(userID uint) ([]models.Chat, error)
	GetLastMessage(chatID uint) (*models.Message, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) FindByUsers(user1ID, user2ID uint) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) Create(chat *models.Chat) error {
	if chat == nil {
		return errors.New("nil chat")
	}
	return r.db.Create(chat).Error
}

func (r *chatRepository) GetUserChats(userID uint) ([]models.Chat, error) {
	var chats []models.Chat
	err := r.db.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *chatRepository) GetLastMessage(chatID uint) (*models.Message, error) {
	var msg models.Message
	err := r.db.Where("chat_id = ?", chatID).Order("created_at DESC").First(&msg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}
