package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Online   bool
	LastSeenAt *time.Time
	Typing     bool
}

type Chat struct {
	gorm.Model
	User1ID uint
	User2ID uint 
}

type Message struct {
	gorm.Model
	ChatID   uint
	SenderID uint
	Type     string
	Text     string
	VoiceURL string
	Edited   bool
	Deleted   bool
	Status   string
	DurationSec int
}
