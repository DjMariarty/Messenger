package services

import (
	"errors"
	"log/slog"

	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/models"
	"github.com/DjMariarty/messenger/internal/repository"
)

var (
	ErrInvalidChatID   = errors.New("chatID cannot be 0")
	ErrInvalidSenderID = errors.New("senderID cannot be 0")
	ErrEmptyMessage    = errors.New("text  cannot be empty")
)

type MessageService interface {
	CreateMessage(req dto.CreateMessageRequest) (*models.Message, error)
	GetMessagesByChatID(chatID uint) ([]models.Message, error)
}

type messageService struct {
	messages repository.MessageRepository
	log      *slog.Logger
}

func NewMessageService(messages repository.MessageRepository, log *slog.Logger) MessageService {
	return &messageService{messages: messages, log: log}
}

func (s *messageService) CreateMessage(req dto.CreateMessageRequest) (*models.Message, error) {
	if req.ChatID == 0 {
		s.log.Warn("service: invalid chatID")
		return nil, ErrInvalidChatID
	}

	if req.SenderID == 0 {
		s.log.Warn("service: invalid senderID")
		return nil, ErrInvalidSenderID
	}

	if req.Text == "" {
		s.log.Warn("service: empty message text")
		return nil, ErrEmptyMessage
	}
	// ---------------------------------------------------
	msg := &models.Message{
		ChatID:   req.ChatID,
		SenderID: req.SenderID,
		Text:     req.Text,
	}

	err := s.messages.Create(msg)
	if err != nil {
		s.log.Error("service: failed to create message", "chat_id", req.ChatID, "sender_id", req.SenderID, "error", err)
		return nil, err
	}
	s.log.Info("service: message created", "message_id", msg.ID, "chat_id", msg.ChatID, "sender_id", msg.SenderID)
	return msg, nil
}

func (s *messageService) GetMessagesByChatID(userID uint) ([]models.Message, error) {
	if userID == 0 {
		s.log.Warn("service: invalid userID (0)")
		return nil, errors.New("chatID cannot be 0")
	}

	messages, err := s.messages.GetMessagesByChatID(userID)
	if err != nil {
		s.log.Error("service: failed to fetch messages", "user_id", userID, "error", err)
		return nil, err
	}

	s.log.Info("service: message fetched", "user_id", userID, "count", len(messages))
	return messages, nil
}
