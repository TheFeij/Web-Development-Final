package models

import "time"

type Channel struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	Owner     uint
}
