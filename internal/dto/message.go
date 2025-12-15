package dto

type CreateMessageRequest struct {
	ChatID   uint   `json:"chat_id"`
	SenderID uint   `json:"sender_id"`
	Text     string `json:"text"`
}
