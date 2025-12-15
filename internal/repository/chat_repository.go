package repository

import (
    "log"

    "github.com/DjMariarty/messenger/internal/models"
    "gorm.io/gorm"
)

type ChatRepository struct {
    db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
    return &ChatRepository{db: db}
}

// Получить все чаты пользователя
func (r *ChatRepository) GetUserChats(userID uint) ([]models.Chat, error) {
    log.Printf("[ChatRepository] GetUserChats userID=%d", userID)

    var chats []models.Chat
    err := r.db.
        Where("user1_id = ? OR user2_id = ?", userID, userID).
        Find(&chats).Error

    if err != nil {
        log.Printf("[ChatRepository] DB error: %v", err)
        return nil, err
    }

    return chats, nil
}

// Получить имя пользователя
func (r *ChatRepository) GetUserNameByID(userID uint) (string, error) {
    log.Printf("[ChatRepository] GetUserNameByID userID=%d", userID)

    var name string
    err := r.db.
        Model(&models.User{}).
        Select("name").
        Where("id = ?", userID).
        Scan(&name).Error

    if err != nil {
        log.Printf("[ChatRepository] DB error: %v", err)
        return "", err
    }

    return name, nil
}

// Получить последнее сообщение чата
func (r *ChatRepository) GetLastMessage(chatID uint) (*models.Message, error) {
    log.Printf("[ChatRepository] GetLastMessage chatID=%d", chatID)

    var msg models.Message
    err := r.db.
        Where("chat_id = ?", chatID).
        Order("created_at DESC").
        First(&msg).Error

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        log.Printf("[ChatRepository] DB error: %v", err)
        return nil, err
    }

    return &msg, nil
}
