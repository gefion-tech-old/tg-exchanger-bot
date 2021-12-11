package models

type Bill struct {
	ID        uint   `json:"id"`
	ChatID    int64  `json:"chat_id"`
	Bill      string `json:"bill"`
	CreatedAt string `json:"created_at"`
}
