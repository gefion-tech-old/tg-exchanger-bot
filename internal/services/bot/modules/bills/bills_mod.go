package bills

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ ModBillsI = (*ModBills)(nil)

type ModBills struct {
	bAPI *tgbotapi.BotAPI
	sAPI api.ApiI
	kbd  keyboards.KeyboardsI
}

type ModBillsI interface {
	/* CallbackQuery обработчики */

	/* Message обработчики */
	MyBills(ctx context.Context, update tgbotapi.Update) error
}

func InitModBills(bAPI *tgbotapi.BotAPI, servAPI api.ApiI, k keyboards.KeyboardsI) ModBillsI {
	return &ModBills{
		bAPI: bAPI,
		sAPI: servAPI,
		kbd:  k,
	}
}
