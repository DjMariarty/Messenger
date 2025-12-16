package services

import (
	"errors"
	"log/slog"

	"github.com/DjMariarty/messenger/internal/auth"
	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/models"
	"github.com/DjMariarty/messenger/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidCredentials = errors.New("неверный email или пароль")
)

type UserService interface {
	RegisterUser(req dto.RegisterRequest) (*models.User, error)

	LoginUser(data dto.LoginRequest) (string, error)

	GetByID(id uint) (*models.User, error)
}

type userService struct {
	users repository.UserRepository
	db    *gorm.DB
	log   *slog.Logger
}

func NewUserService(
	db *gorm.DB,
	users repository.UserRepository,
	log *slog.Logger,

) UserService {
	return &userService{
		db:    db,
		users: users,
		log:   log,
	}
}

func (s *userService) RegisterUser(req dto.RegisterRequest) (*models.User, error) {
	s.log.Info("user service: register started", slog.String("email", req.Email))
	if err := s.validateUserRegister(req); err != nil {
		s.log.Warn("user service: register validation failed",
			slog.String("email", req.Email),
			slog.String("reason", err.Error()),
		)
		return nil, err
	}

	var createdUser models.User

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var existing models.User
		if err := tx.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			s.log.Warn("user service: register failed - email already exists",
				slog.String("email", req.Email),
			)
			return errors.New("пользователь с таким email уже существует")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("user service: register failed - db error on email check",
				slog.String("email", req.Email),
				slog.Any("error", err),
			)
			return err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.log.Error("user service: register failed - bcrypt error",
				slog.String("email", req.Email),
				slog.Any("error", err),
			)
			return err
		}

		createdUser = models.User{
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: string(hash),
		}

		if err := tx.Create(&createdUser).Error; err != nil {
			s.log.Error("user service: register failed - create user error",
				slog.String("email", req.Email),
				slog.Any("error", err),
			)
			return err
		}

		return nil
	})

	if err != nil {

		return nil, err
	}

	s.log.Info("user service: register success",
		slog.Uint64("user_id", uint64(createdUser.ID)),
		slog.String("email", createdUser.Email),
	)

	return &createdUser, nil
}

func (s *userService) LoginUser(req dto.LoginRequest) (string, error) {
	s.log.Info("user service: login started", slog.String("email", req.Email))

	user, err := s.users.GetByEmail(req.Email)
	if err != nil {
		s.log.Warn("user service: login failed - invalid credentials",
			slog.String("email", req.Email),
		)
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.log.Warn("user service: login failed - invalid credentials",
			slog.String("email", req.Email),
			slog.Uint64("user_id", uint64(user.ID)),
		)
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		s.log.Error("user service: login failed - token generation error",
			slog.Uint64("user_id", uint64(user.ID)),
			slog.Any("error", err),
		)
		return "", err
	}

	s.log.Info("user service: login success",
		slog.Uint64("user_id", uint64(user.ID)),
	)

	return token, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	s.log.Info("user service: get by id started",
		slog.Uint64("user_id", uint64(id)),
	)

	user, err := s.users.GetByID(id)
	if err != nil {
		s.log.Error("user service: get by id failed",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.log.Info("user service: get by id success",
		slog.Uint64("user_id", uint64(id)),
	)

	return user, nil
}

func (s *userService) validateUserRegister(req dto.RegisterRequest) error {
	if req.Name == "" {
		return errors.New("имя не может быть пустым")
	}
	if req.Email == "" {
		return errors.New("эмайл не может быть пустым")
	}
	if req.Password == "" {
		return errors.New("пароль не может быть пустым")
	}
	return nil
}
