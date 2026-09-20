package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	From      string    `json:"from"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}
