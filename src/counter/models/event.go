package models

type CounterEvent struct {
	EventType   string `json:"eventType"`
	DialogID    string `json:"dialogID"`
	RecipientID string `json:"recipientID"`
	MessageID   int    `json:"messageID"`
}
