package transport

import (
	"net/http"

	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	users services.UserService
}

func NewUserHandler(users services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.users.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.UserResponse{
		ID: user.ID, Name: user.Name, Email: user.Email,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.users.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	user, err := h.users.GetByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
