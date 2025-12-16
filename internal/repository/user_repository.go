package repository

import (
	"errors"
	"log/slog"

	"github.com/DjMariarty/messenger/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type gormUserRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewUserRepository(db *gorm.DB, log *slog.Logger) UserRepository {
	return &gormUserRepository{
		db:  db,
		log: log,
	}
}

func (r *gormUserRepository) Create(user *models.User) error {
	if user == nil {
		r.log.Warn(
			"user repository: create called with nil user",
		)
		return errors.New("пустой пользователь")
	}

	if err := r.db.Create(user).Error; err != nil {
		r.log.Error(
			"user repository: failed to create user",
			slog.String("email", user.Email),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (r *gormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warn(
				"user repository: user not found by id",
				slog.Uint64("user_id", uint64(id)),
			)
		} else {
			r.log.Error(
				"user repository: failed to get user by id",
				slog.Uint64("user_id", uint64(id)),
				slog.Any("error", err),
			)
		}
		return nil, err
	}

	return &user, nil
}

func (r *gormUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warn(
				"user repository: user not found by email",
				slog.String("email", email),
			)
		} else {
			r.log.Error(
				"user repository: failed to get user by email",
				slog.String("email", email),
				slog.Any("error", err),
			)
		}
		return nil, err
	}

	return &user, nil
}
