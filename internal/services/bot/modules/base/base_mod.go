package base

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ModBase struct {
	bAPI  *tgbotapi.BotAPI
	sAPI  api.ApiI
	redis redisstore.RedisStoreI
	kbd   keyboards.KeyboardsI
}

type ModBaseI interface {
	/* ReplyKeyboardMarkup обработчики*/
	SupportRequest(ctx context.Context, update tgbotapi.Update) error
	AboutBot(ctx context.Context, update tgbotapi.Update) error
	Operators(ctx context.Context, update tgbotapi.Update) error
}

func InitModBase(bAPI *tgbotapi.BotAPI, servAPI api.ApiI, redis redisstore.RedisStoreI, k keyboards.KeyboardsI) ModBaseI {
	return &ModBase{
		bAPI:  bAPI,
		sAPI:  servAPI,
		redis: redis,
		kbd:   k,
	}
}
