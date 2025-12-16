package services

import (
	"errors"
	"sort"
	"time"

	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/models"
	"github.com/DjMariarty/messenger/internal/repository"
	"gorm.io/gorm"
)

type ChatService interface {
	CreateChat(userID uint, req dto.CreateChatRequest) (*models.Chat, error)
	GetChats(userID uint) ([]dto.ChatResponse, error)
}

type chatService struct {
	db    *gorm.DB
	chats repository.ChatRepository
}

func NewChatService(db *gorm.DB, chats repository.ChatRepository) ChatService {
	return &chatService{db: db, chats: chats}
}

func (s *chatService) CreateChat(userID uint, req dto.CreateChatRequest) (*models.Chat, error) {
	if userID == 0 || req.PartnerID == 0 {
		return nil, errors.New("invalid user id")
	}
	if userID == req.PartnerID {
		return nil, errors.New("cannot create chat with yourself")
	}


	u1, u2 := userID, req.PartnerID
	if u1 > u2 {
		u1, u2 = u2, u1
	}


	existing, err := s.chats.FindByUsers(u1, u2)
	if err == nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}


	chat := models.Chat{User1ID: u1, User2ID: u2}
	if err := s.chats.Create(&chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (s *chatService) GetChats(userID uint) ([]dto.ChatResponse, error) {
	chats, err := s.chats.GetUserChats(userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ChatResponse, 0, len(chats))

	for _, ch := range chats {
		lastMsg, err := s.chats.GetLastMessage(ch.ID)
		if err != nil {
			return nil, err
		}

		var lastText string
		var lastTime *time.Time
		if lastMsg != nil {
			lastText = lastMsg.Text
			t := lastMsg.CreatedAt
			lastTime = &t
		}

		res = append(res, dto.ChatResponse{
			ChatID:          ch.ID,
			LastMessage:     lastText,
			LastMessageTime: lastTime,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		ti := res[i].LastMessageTime
		tj := res[j].LastMessageTime
		if ti == nil && tj == nil {
			return false
		}
		if ti == nil {
			return false
		}
		if tj == nil {
			return true
		}
		return ti.After(*tj)
	})

	return res, nil
}
