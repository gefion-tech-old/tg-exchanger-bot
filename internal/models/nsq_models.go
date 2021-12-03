package models

type MessageEvent struct {
	Message struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"message"`
	To struct {
		ChatID   int64  `json:"chat_id"`
		Username string `json:"username"`
	} `json:"to"`
	CreatedAt string `json:"created_at"`
}
