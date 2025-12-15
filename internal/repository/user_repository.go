package repository

import (
	"errors"

	"github.com/DjMariarty/messenger/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error

	GetByID(id uint) (*models.User, error)

	GetByEmail(email string) (*models.User, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *models.User) error {
	if user == nil {
		return errors.New("пустой пользователь")
	}
	return r.db.Create(user).Error
}

func (r *gormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
