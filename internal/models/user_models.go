package models

type UserReq struct {
	ChatID   int
	Username string
}

type UserAction struct {
	ActionType int
	Step       int
	MetaData   map[string]interface{}
	User       UserReq
}
