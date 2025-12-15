package services

import (
	"errors"

	"github.com/DjMariarty/messenger/internal/auth"
	"github.com/DjMariarty/messenger/internal/dto"
	"github.com/DjMariarty/messenger/internal/models"
	"github.com/DjMariarty/messenger/internal/repository"
	"golang.org/x/crypto/bcrypt"
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
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (r *userService) validateUserRegister(req dto.RegisterRequest) error {
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

func (s *userService) RegisterUser(req dto.RegisterRequest) (*models.User, error) {
	if err := s.validateUserRegister(req); err != nil {
		return nil, err
	}

	if _, err := s.users.GetByEmail(req.Email); err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) LoginUser(req dto.LoginRequest) (string, error) {
	user, err := s.users.GetByEmail(req.Email)
	if err != nil {
		return "", ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("неверный email или пароль")

	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.users.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, err
}
