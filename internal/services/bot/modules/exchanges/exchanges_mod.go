package exchanges

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ ModExchangesI = (*ModExchanges)(nil)

type ModExchanges struct {
	bAPI *tgbotapi.BotAPI
	sAPI api.ApiI
	kbd  keyboards.KeyboardsI
}

type ModExchangesI interface {
	/* CallbackQuery обработчики */

	// @CallbackQuery BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE
	ReceiveAsResultOfExchange(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error

	// @CallbackQuery BOT__CQ__EX__REQ_AMOUNT
	ReqAmount(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error

	/*  Универсальные обработчики */

	// @Button BOT__BTN__BASE__NEW_EXCHANGE
	// @CallbackQuery BOT__CQ__EX__COINS_TO_EXCHAGE
	NewExchange(ctx context.Context, update tgbotapi.Update) error
}

func InitModExchanges(bAPI *tgbotapi.BotAPI, servAPI api.ApiI, k keyboards.KeyboardsI) ModExchangesI {
	return &ModExchanges{
		bAPI: bAPI,
		sAPI: servAPI,
		kbd:  k,
	}
}
