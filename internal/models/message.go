package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	ChatID   uint   `json:"chat_id"`
	SenderID uint   `json:"sender_id"`
	Text     string `json:"text"`
}
