package models

import "time"

type Message struct {
	From      string    `json:"from"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
