package models

import "gorm.io/gorm"

type Chat struct {
    gorm.Model

    User1ID uint `json:"user1_id"`
    User2ID uint `json:"user2_id"`
	    User1   User `gorm:"foreignKey:User1ID"`  
    User2   User `gorm:"foreignKey:User2ID"`
}
