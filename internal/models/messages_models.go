package models

type Message struct {
	ID          uint   `json:"id"`
	Connector   string `json:"connector"`
	MessageText string `json:"message_text"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
