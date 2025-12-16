package dto

import "time"

type CreateChatRequest struct {
	PartnerID uint `json:"partner_id"`
}

type CreateChatResponse struct {
	ChatID  uint `json:"chat_id"`
	User1ID uint `json:"user1_id"`
	User2ID uint `json:"user2_id"`
}

type ChatResponse struct {
	ChatID          uint       `json:"chat_id"`
	LastMessage     string     `json:"last_message"`
	LastMessageTime *time.Time `json:"last_message_time"`
}
