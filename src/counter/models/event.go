package models

type CounterEvent struct {
	EventType   string `json:"eventType"`
	DialogID    int    `json:"dialogID"`
	RecipientID int    `json:"recipientID"`
	MessageID   int    `json:"messageID"`
}
