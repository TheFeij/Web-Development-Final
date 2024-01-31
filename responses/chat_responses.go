package responses

import (
	"time"
)

type Chat struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IsDead    bool      `json:"is_dead"`
}

type ChatsList struct {
	Chats []Chat `json:"chats"`
}

type ChatContent struct {
	Chat     Chat          `json:"chat"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	ID                uint      `json:"id"`
	ChatID            uint      `json:"chat_id"`
	SourceSenderID    uint      `json:"source_sender_id"`
	OriginalMessageID uint      `json:"original_message_id"`
	SenderID          uint      `json:"sender_id"`
	ReceiverID        uint      `json:"receiver_id"`
	Content           string    `json:"content"`
	CreatedAt         time.Time `json:"created_at"`
}
