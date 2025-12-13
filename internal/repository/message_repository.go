package repository

import (
	"errors"
	"log/slog"

	"github.com/DjMariarty/messenger/internal/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *models.Message) error
	GetMessagesByChatID(chatID uint) ([]models.Message, error)
}

var ErrMessageNil = errors.New("message nil")

type gormMessageRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewMessageRepository(db *gorm.DB, log *slog.Logger) MessageRepository {
	return &gormMessageRepository{db: db, log: log}
}

func (r *gormMessageRepository) Create(message *models.Message) error {
	if message == nil {
		r.log.Error("create: message is nil")
		return ErrMessageNil
	}

	r.log.Debug("creating message", "chat_id", message.ChatID, "sender_id", message.SenderID)

	if err := r.db.Create(message).Error; err != nil {
		r.log.Error("create failed", "chat_id", message.ChatID, "sender_id", message.SenderID, "error", err)
		return err
	}

	return nil

}

func (r *gormMessageRepository) GetMessagesByChatID(userID uint) ([]models.Message, error) {
	r.log.Debug("fetch messages by chat", "user_id", userID)

	var messages []models.Message
	err := r.db.Model(&models.Message{}).Where("user_id = ?", userID).Order("created_at desc").Find(&messages).Error
	if err != nil {
		r.log.Error("fetch messages failed", "user_id", userID, "error", err)
		return nil, err
	}
	return messages, err
}
