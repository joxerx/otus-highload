package models

type Message struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type MessageRequest struct {
	Text string `json:"text"`
}
