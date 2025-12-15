package services

import (
    "errors"
    "log"
    "sort"
    "time"

    "github.com/DjMariarty/messenger/internal/models"
    "github.com/DjMariarty/messenger/internal/repository"
)

var (
    ErrInvalidUserID = errors.New("invalid user id")
)

type ChatPreview struct {
    ChatID uint `json:"chat_id"`

    Partner struct {
        ID   uint   `json:"id"`
        Name string `json:"name"`
    } `json:"partner"`

    LastMessage     string    `json:"last_message,omitempty"`
    LastMessageTime time.Time `json:"last_message_time,omitempty"`
}

type ChatService struct {
    repo *repository.ChatRepository
}

func NewChatService(repo *repository.ChatRepository) *ChatService {
    return &ChatService{repo: repo}
}

func (s *ChatService) GetUserChatList(userID uint) ([]ChatPreview, error) {
    log.Printf("[ChatService] GetUserChatList userID=%d", userID)

    if userID == 0 {
        return nil, ErrInvalidUserID
    }

    chats, err := s.repo.GetUserChats(userID)
    if err != nil {
        return nil, err
    }

    previews := make([]ChatPreview, 0, len(chats))

    for _, chat := range chats {
        preview, err := s.buildPreview(chat, userID)
        if err != nil {
            log.Printf("[ChatService] buildPreview error: %v", err)
            continue
        }
        previews = append(previews, preview)
    }

    sort.Slice(previews, func(i, j int) bool {
        return previews[i].LastMessageTime.After(previews[j].LastMessageTime)
    })

    return previews, nil
}

func (s *ChatService) buildPreview(chat models.Chat, currentUser uint) (ChatPreview, error) {
    var preview ChatPreview
    preview.ChatID = chat.ID

    partnerID := chat.User1ID
    if chat.User1ID == currentUser {
        partnerID = chat.User2ID
    }

    name, err := s.repo.GetUserNameByID(partnerID)
    if err != nil {
        return preview, err
    }

    preview.Partner.ID = partnerID
    preview.Partner.Name = name

    lastMsg, err := s.repo.GetLastMessage(chat.ID)
    if err != nil {
        return preview, err
    }

    if lastMsg != nil {
        preview.LastMessage = lastMsg.Text
        preview.LastMessageTime = lastMsg.CreatedAt
    }

    return preview, nil
}
