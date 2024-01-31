package models

import "time"

type ChatMessage struct {
	ID                uint
	ChatID            uint
	SourceSenderID    uint
	OriginalMessageID uint
	SenderID          uint
	ReceiverID        uint
	Content           string
	CreatedAt         time.Time
}
