package models

type Exchanger struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UrlToParse string `json:"url"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
