package models

import "time"

type Post struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
