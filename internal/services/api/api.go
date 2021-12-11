package api

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
)

type Api struct {
	config *config.ApiConfig

	userReq UserRequestsI
	msgReq  MessageRequestsI
	billReq BillRequestsI
}

type ApiI interface {
	User() UserRequestsI
	Message() MessageRequestsI
	Bill() BillRequestsI
}

func Init(c *config.ApiConfig) ApiI {
	return &Api{
		config: c,
	}
}

func (api *Api) Bill() BillRequestsI {
	if api.billReq != nil {
		return api.billReq
	}

	api.billReq = InitBillRequests(api.config.Url)
	return api.billReq
}

func (api *Api) Message() MessageRequestsI {
	if api.msgReq != nil {
		return api.msgReq
	}

	api.msgReq = InitMessageRequests(api.config.Url)
	return api.msgReq
}

func (api *Api) User() UserRequestsI {
	if api.userReq != nil {
		return api.userReq
	}

	api.userReq = InitUserRequests(api.config.Url)
	return api.userReq
}
