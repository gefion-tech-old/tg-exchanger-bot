package api

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
)

type Api struct {
	config  *config.ApiConfig
	userReq *UserRequests
	msgReq  *MessageRequests
}

type ApiI interface {
	User() UserRequestsI
	Message() MessageRequestsI
}

func Init(c *config.ApiConfig) ApiI {
	return &Api{
		config: c,
	}
}

func (api *Api) Message() MessageRequestsI {
	if api.msgReq != nil {
		return api.msgReq
	}

	api.msgReq = &MessageRequests{
		url: api.config.Url,
	}
	return api.msgReq
}

func (api *Api) User() UserRequestsI {
	if api.userReq != nil {
		return api.userReq
	}

	api.userReq = &UserRequests{
		url: api.config.Url,
	}
	return api.userReq
}
