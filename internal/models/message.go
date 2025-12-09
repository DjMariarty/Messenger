package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	Chat_id   uint
	Sender_id uint
	Text      string
}
