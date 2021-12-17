package bills

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ ModBillsI = (*ModBills)(nil)

type ModBills struct {
	bAPI *tgbotapi.BotAPI
	sAPI api.ApiI
	cnf  *config.BotConfig

	redis redisstore.RedisStoreI
	kbd   keyboards.KeyboardsI
}

type ModBillsI interface {
	/* CallbackQuery обработчики */
	AddNewBillStepOne(ctx context.Context, update tgbotapi.Update) error
	AddNewBillStepTwo(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error
	AddNewBillStepThreeCorrect(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error
	AddNewBillStepThreeInCorrect(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error
	AddNewBillStepFour(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error

	/* Message обработчики */
	MyBills(ctx context.Context, update tgbotapi.Update) error
}

func InitModBills(bAPI *tgbotapi.BotAPI, servAPI api.ApiI, redis redisstore.RedisStoreI, k keyboards.KeyboardsI, cnf *config.BotConfig) ModBillsI {
	return &ModBills{
		bAPI:  bAPI,
		sAPI:  servAPI,
		redis: redis,
		kbd:   k,
		cnf:   cnf,
	}
}
