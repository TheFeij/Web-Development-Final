package models

import "time"

type Group struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	OwnerID   uint
}
