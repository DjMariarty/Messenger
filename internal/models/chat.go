package models

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	User1ID   uint `gorm:"not null"`
	User2ID   uint `gorm:"not null"`
	CreatedAt time.Time

	User1 User `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE"`
	User2 User `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE"`

	Messages []Message `gorm:"foreignKey:ChatID"`
}
