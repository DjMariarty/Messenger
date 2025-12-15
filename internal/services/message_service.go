 package services

// import (
// 	"errors"

// 	"github.com/DjMariarty/messenger/internal/models"
// 	"github.com/DjMariarty/messenger/internal/repository"
// )

// type MessageService interface {
// 	CreateMessage(req models.CreateMessageRequest) (*models.Message, error)
// }

// type messageService struct {
// 	messages repository.MessageRepository
// }

// func NewMessageService(messages repository.MessageRepository) MessageService {
// 	return &messageService{messages: messages}
// }

// func (s *messageService) CreateMessage(req models.CreateMessageRequest) (*models.Message, error) {
// 	if req.Text == "" {
// 		return nil, errors.New("text cannot be empty")
// 	}

// 	msg := &models.Message{
// 		ChatID:   req.ChatID,
// 		SenderID: req.SenderID,
// 		Text:     req.Text,
// 	}

// 	err := s.messages.Create(msg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return msg, nil
// }
