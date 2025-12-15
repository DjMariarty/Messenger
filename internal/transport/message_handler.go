package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/services"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	service services.MessageService
	log     *slog.Logger
}

func NewMessageHandler(service services.MessageService, log *slog.Logger) *MessageHandler {
	return &MessageHandler{service: service, log: log}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	// логи на входе

	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("handler: invalid JSON in CreateMessage", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.service.CreateMessage(req)
	if err != nil {
		h.log.Error("handler: failed to create message", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("handler: creating message",
		slog.Uint64("chat_id", uint64(req.ChatID)),
		slog.Uint64("sender_id", uint64(req.SenderID)),
	)
	c.JSON(http.StatusCreated, msg)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userIDParam := c.Param("userID")
	if userIDParam == "" {
		h.log.Warn("handler: missing chat_id in path")
		c.JSON(http.StatusBadRequest, gin.H{"error": "chatID id required"})
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		h.log.Warn("handler: invalid user_id format", slog.String("user_id", userIDParam))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chatID"})
		return
	}

	messages, err := h.service.GetMessagesByChatID(uint(userID))
	if err != nil {
		h.log.Error("handler: failed to get messages",
			slog.Uint64("user_id", userID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("handler: messages fetched successfully",
		slog.Uint64("user_id", userID),
		slog.Int("count", len(messages)),
	)

	c.JSON(http.StatusOK, messages)

}
