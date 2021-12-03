package btns

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Btns struct {
	botAPI *tgbotapi.BotAPI
	sAPI   api.ApiI
}

type BtnsI interface{}

func Init(bAPI *tgbotapi.BotAPI, sAPI api.ApiI) BtnsI {
	return &Btns{
		botAPI: bAPI,
		sAPI:   sAPI,
	}
}
