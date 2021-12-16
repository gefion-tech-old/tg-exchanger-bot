package models

type UserAction struct {
	ActionType int
	Step       int
	MetaData   map[string]interface{}
	User       struct {
		ChatID   int
		Username string
	}
}
