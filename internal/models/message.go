package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatID   uint   `json:"chat_id" gorm:"not null;index"`
	SenderID uint   `json:"sender_id" gorm:"not null;index"`
	Text     string `json:"text" gorm:"not null"`
}
