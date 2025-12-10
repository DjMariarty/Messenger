package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model

	User1ID uint `json:"user1_id"`
	User2ID uint `json:"user2_id"`
}
