package repository

import (
	"github.com/DjMariarty/messenger/internal/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(Message *models.Message) error
}

type gormMessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &gormMessageRepository{db: db}
}

func (r *gormMessageRepository) Create(message *models.Message) error {
	if message == nil {
		return nil
	}
	return r.db.Create(message).Error
}
