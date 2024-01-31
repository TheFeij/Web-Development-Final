package models

import "time"

type ChannelPost struct {
	ID                uint
	ChannelID         uint
	SourceSenderID    uint
	OriginalMessageID uint
	SenderID          uint
	Content           string
	CreatedAt         time.Time
}
