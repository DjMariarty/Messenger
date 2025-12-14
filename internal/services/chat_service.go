package service

import (
	"log"
	"sort"
	"time"

	"github.com/DjMariarty/messenger/internal/models"
	"github.com/DjMariarty/messenger/internal/repository"
	"gorm.io/gorm"
)

type ChatPreview struct {
	ChatID          uint      `json:"chat_id"`
	Partner         struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"partner"`
	LastMessage     string    `json:"last_message,omitempty"`
	LastMessageTime time.Time `json:"last_message_time,omitempty"`
}

type ChatService struct {
	db       *gorm.DB
	chatRepo *repository.ChatRepository
}

func NewChatService(db *gorm.DB, chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{
		db:       db,
		chatRepo: chatRepo,
	}
}

func (s *ChatService) GetUserChatList(userID uint) ([]ChatPreview, error) {
	if userID == 0 {
		log.Println("[ChatService] GetUserChatList: invalid userID = 0")
		return nil, gorm.ErrInvalidData
	}

	chats, err := s.chatRepo.GetUserChats(userID)
	if err != nil {
		log.Printf("[ChatService] GetUserChatList: failed to get chats for userID=%d, err=%v\n", userID, err)
		return nil, err
	}

	if len(chats) == 0 {
		log.Printf("[ChatService] GetUserChatList: no chats found for userID=%d\n", userID)
		return []ChatPreview{}, nil
	}

	previews := make([]ChatPreview, 0, len(chats))
	for _, chat := range chats {
		preview, err := s.createChatPreview(&chat, userID)
		if err != nil {
			log.Printf("[ChatService] createChatPreview: skipping chatID=%d, err=%v\n", chat.ID, err)
			continue
		}
		previews = append(previews, *preview)
	}

	sortChatsByLastMessage(previews)
	log.Printf("[ChatService] GetUserChatList: returning %d previews for userID=%d\n", len(previews), userID)
	return previews, nil
}

func (s *ChatService) createChatPreview(chat *models.Chat, currentUserID uint) (*ChatPreview, error) {
	var partnerID uint
	if chat.User1ID == currentUserID {
		partnerID = chat.User2ID
	} else {
		partnerID = chat.User1ID
	}

	if partnerID == 0 {
		log.Printf("[ChatService] createChatPreview: invalid partnerID for chatID=%d\n", chat.ID)
		return nil, gorm.ErrInvalidData
	}

	var partnerName string
	err := s.db.Model(&models.User{}).
		Select("name").
		Where("id = ?", partnerID).
		Scan(&partnerName).Error
	if err != nil {
		log.Printf("[ChatService] createChatPreview: failed to get partner name for userID=%d, err=%v\n", partnerID, err)
		return nil, err
	}

	lastMessage, err := s.chatRepo.GetLastMessage(chat.ID)
	if err != nil {
		log.Printf("[ChatService] createChatPreview: failed to get last message for chatID=%d, err=%v\n", chat.ID, err)
		return nil, err
	}

	preview := &ChatPreview{
		ChatID: chat.ID,
	}
	preview.Partner.ID = partnerID
	preview.Partner.Name = partnerName
	if lastMessage != nil {
		preview.LastMessage = lastMessage.Text
		preview.LastMessageTime = lastMessage.CreatedAt
	}

	return preview, nil
}

func sortChatsByLastMessage(chats []ChatPreview) {
	sort.Slice(chats, func(i, j int) bool {
		ti := chats[i].LastMessageTime
		tj := chats[j].LastMessageTime

		if ti.IsZero() && tj.IsZero() {
			return false
		}
		if ti.IsZero() {
			return false
		}
		if tj.IsZero() {
			return true
		}

		return ti.After(tj)
	})
}
