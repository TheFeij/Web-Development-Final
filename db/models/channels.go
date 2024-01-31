package models

import "time"

type Channel struct {
	ID        uint
	CreatedAt time.Time
	Owner     uint
}
