package api

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/config"
	"github.com/sirupsen/logrus"
)

type Api struct {
	config  *config.ApiConfig
	bConfig *config.BotConfig
	logger  *logrus.Logger

	userReq         UserRequestsI
	msgReq          MessageRequestsI
	billReq         BillRequestsI
	notificationReq NotificationRequestsI
	exchangerReq    ExchangerRequestsI

	tgReq TelegramRequestsI
}

type ApiI interface {
	User() UserRequestsI
	Message() MessageRequestsI
	Bill() BillRequestsI
	Notification() NotificationRequestsI
	Telegram() TelegramRequestsI
	Exchanger() ExchangerRequestsI
}

func Init(c *config.ApiConfig, bC *config.BotConfig, l *logrus.Logger) ApiI {
	return &Api{
		config:  c,
		bConfig: bC,
		logger:  l,
	}
}

func (api *Api) Exchanger() ExchangerRequestsI {
	if api.exchangerReq != nil {
		return api.exchangerReq
	}

	api.exchangerReq = InitExchangerRequests(api.config.Url, api.logger)
	return api.exchangerReq
}

func (api *Api) Telegram() TelegramRequestsI {
	if api.tgReq != nil {
		return api.tgReq
	}

	api.tgReq = InitTelegramRequests("https://api.telegram.org/", api.bConfig.Token, api.logger)
	return api.tgReq
}

func (api *Api) Notification() NotificationRequestsI {
	if api.notificationReq != nil {
		return api.notificationReq
	}

	api.notificationReq = InitNotificationRequests(api.config.Url, api.logger)
	return api.notificationReq
}

func (api *Api) Bill() BillRequestsI {
	if api.billReq != nil {
		return api.billReq
	}

	api.billReq = InitBillRequests(api.config.Url, api.logger)
	return api.billReq
}

func (api *Api) Message() MessageRequestsI {
	if api.msgReq != nil {
		return api.msgReq
	}

	api.msgReq = InitMessageRequests(api.config.Url, api.logger)
	return api.msgReq
}

func (api *Api) User() UserRequestsI {
	if api.userReq != nil {
		return api.userReq
	}

	api.userReq = InitUserRequests(api.config.Url, api.logger)
	return api.userReq
}
