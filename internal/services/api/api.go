package api

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
)

type Api struct {
	config  *config.ApiConfig
	userReq *UserRequests
}

type ApiI interface {
	/*
		Сботка методов для работы с пользователскими данными и аккаунтом
	*/
	User() UserRequestsI
}

func Init(c *config.ApiConfig) ApiI {
	return &Api{
		config: c,
	}
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
