package exchanges

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

var _ ModExchangesI = (*ModExchanges)(nil)

type ModExchanges struct {
	bAPI   *tgbotapi.BotAPI
	sAPI   api.ApiI
	redis  redisstore.RedisStoreI
	kbd    keyboards.KeyboardsI
	logger *logrus.Logger
}

type ModExchangesI interface {
	/* Action обработчики */

	CreateLinkForPayment(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error

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

func InitModExchanges(bAPI *tgbotapi.BotAPI, servAPI api.ApiI, redis redisstore.RedisStoreI, k keyboards.KeyboardsI, l *logrus.Logger) ModExchangesI {
	return &ModExchanges{
		bAPI:   bAPI,
		sAPI:   servAPI,
		redis:  redis,
		kbd:    k,
		logger: l,
	}
}
