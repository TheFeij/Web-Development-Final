package models

import "time"

type GroupMessage struct {
	ID                uint
	GroupID           uint
	SourceSenderID    uint
	OriginalMessageID uint
	SenderID          uint
	Content           string
	CreatedAt         time.Time
}
