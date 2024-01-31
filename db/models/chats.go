package models

import (
	"time"
)

type Chat struct {
	ID        uint
	CreatedAt time.Time
	IsDead    bool
}
