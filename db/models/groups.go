package models

import "time"

type Groups struct {
	ID        uint
	CreatedAt time.Time
	Owner     uint
}
