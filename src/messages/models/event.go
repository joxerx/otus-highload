package models

type ErrorEvent struct {
	EventType string `json:"eventType"`
	MessageID int    `json:"messageId"`
}
