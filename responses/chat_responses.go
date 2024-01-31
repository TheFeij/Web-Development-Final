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
