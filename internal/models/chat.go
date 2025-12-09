package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model

	User1_id uint
	User2_id uint
}
