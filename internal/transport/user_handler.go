package transport

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	users services.UserService
	log   *slog.Logger
}

func NewUserHandler(users services.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		users: users,
		log:   log,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("user handler: register bind json failed",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	h.log.Info("user handler: register request received",
		slog.String("email", req.Email),
	)

	user, err := h.users.RegisterUser(req)
	if err != nil {

		if err.Error() == "пользователь с таким email уже существует" {
			h.log.Warn("user handler: register conflict - email exists",
				slog.String("email", req.Email),
			)
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.log.Warn("user handler: register failed - record not found",
				slog.String("email", req.Email),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.log.Error("user handler: register failed",
			slog.String("email", req.Email),
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("user handler: register success",
		slog.Uint64("user_id", uint64(user.ID)),
		slog.String("email", user.Email),
	)

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("user handler: login bind json failed",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	h.log.Info("user handler: login request received",
		slog.String("email", req.Email),
	)

	token, err := h.users.LoginUser(req)
	if err != nil {

		h.log.Warn("user handler: login failed",
			slog.String("email", req.Email),
			slog.Any("error", err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.log.Info("user handler: login success",
		slog.String("email", req.Email),
	)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	h.log.Info("user handler: me request received",
		slog.Uint64("user_id", uint64(userID)),
	)

	user, err := h.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.log.Warn("user handler: me user not found",
				slog.Uint64("user_id", uint64(userID)),
			)
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		h.log.Error("user handler: me failed",
			slog.Uint64("user_id", uint64(userID)),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.log.Info("user handler: me success",
		slog.Uint64("user_id", uint64(userID)),
	)

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
