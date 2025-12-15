package transport 
import (
	"log"
	"net/http"
	"strconv"

	"github.com/DjMariarty/messenger/internal/services"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
    service *services.ChatService
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
    return &ChatHandler{service: service}
}

func (h *ChatHandler) GetChats(c *gin.Context) {
    userIDValue, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    userID, err := parseUserID(userIDValue)
    if err != nil || userID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    log.Printf("[ChatHandler] GetChats userID=%d", userID)

    chats, err := h.service.GetUserChatList(userID)
    if err != nil {
        log.Printf("[ChatHandler] service error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load chats"})
        return
    }

    c.JSON(http.StatusOK, chats)
}

func parseUserID(v interface{}) (uint, error) {
    switch val := v.(type) {
    case uint:
        return val, nil
    case int:
        return uint(val), nil
    case int64:
        return uint(val), nil
    case float64:
        return uint(val), nil
    case string:
        parsed, err := strconv.ParseUint(val, 10, 32)
        return uint(parsed), err
    default:
        return 0, strconv.ErrSyntax
    }
}
