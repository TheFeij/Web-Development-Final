package responses

import "time"

type Group struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Owner     uint      `json:"owner"`
}
