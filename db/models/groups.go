package models

import "time"

type Groups struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	Owner     uint
}
