package transport

import (
	"net/http"

	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/services"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chats services.ChatService
}

func NewChatHandler(chats services.ChatService) *ChatHandler {
	return &ChatHandler{chats: chats}
}

// POST /chats
func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req dto.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	chat, err := h.chats.CreateChat(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateChatResponse{
		ChatID:  chat.ID,
		User1ID: chat.User1ID,
		User2ID: chat.User2ID,
	})
}

// GET /chats
func (h *ChatHandler) GetChats(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	list, err := h.chats.GetChats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, list)
}
