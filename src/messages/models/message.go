package models

import "time"

type User struct {
	ID          string    `json:"id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_name"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"created_at"`
}
