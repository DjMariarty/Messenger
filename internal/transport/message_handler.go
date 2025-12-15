package transport

// import (
// 	"net/http"

// 	"github.com/DjMariarty/messenger/internal/models"
// 	"github.com/DjMariarty/messenger/internal/services"
// 	"github.com/gin-gonic/gin"
// )

// type MessageHandler struct {
// 	service services.MessageService
// }

// func NewMessageHandler(service services.MessageService) *MessageHandler {
// 	return &MessageHandler{service: service}
// }

// func (h *MessageHandler) CreateMessage(c *gin.Context) {
// 	var req models.CreateMessageRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	msg, err := h.service.CreateMessage(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, msg)
// }
