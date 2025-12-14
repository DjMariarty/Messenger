package handler

import (
	"log"
	"net/http"
	"strconv"

	service "github.com/DjMariarty/messenger/internal/services"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		log.Println("[ChatHandler] GetChats: user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	userID, err := parseUserID(userIDValue)
	if err != nil {
		log.Printf("[ChatHandler] GetChats: invalid user_id value=%v, err=%v\n", userIDValue, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	chats, err := h.chatService.GetUserChatList(userID)
	if err != nil {
		log.Printf("[ChatHandler] GetChats: failed to get chat list for userID=%d, err=%v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список чатов"})
		return
	}

	log.Printf("[ChatHandler] GetChats: returning %d chats for userID=%d\n", len(chats), userID)
	c.JSON(http.StatusOK, chats)
}

func parseUserID(userIDValue interface{}) (uint, error) {
	switch v := userIDValue.(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		parsed, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(parsed), nil
	default:
		return 0, strconv.ErrSyntax
	}
}
