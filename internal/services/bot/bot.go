package bot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/core/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/commands"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
	sAPI   api.ApiI
	cnf    *config.BotConfig
	logger *logrus.Logger
	redis  redisstore.RedisStoreI
	m      modules.BotModulesI
	cmd    commands.CommandsI
	kbd    keyboards.KeyboardsI
}

type BotI interface {
	// Метод слушатель входящих сообщений из очереди
	HandleNsqEvent(consumer *nsq.Consumer, cnf *config.NsqConfig) error
	//Метод слушатель входящих событий в telegram
	HandleBotEvent(ctx context.Context) error
	// Коннектор всех nsq потребителей
	ConnectNsqConsumers(bConsumers *nsqstore.BotConsumers)
}

func Init(bAPI *tgbotapi.BotAPI, sAPI api.ApiI, redis redisstore.RedisStoreI, cnf *config.BotConfig, l *logrus.Logger) BotI {
	kb := keyboards.InitKeyboards()
	mod := modules.InitBotModules(bAPI, kb, redis, sAPI, cnf, l)
	cmd := commands.InitCommands(bAPI, kb, sAPI, l)

	return &Bot{
		botAPI: bAPI,
		sAPI:   sAPI,
		cnf:    cnf,
		logger: l,
		redis:  redis,
		cmd:    cmd,
		kbd:    kb,
		m:      mod,
	}
}

func (bot *Bot) ConnectNsqConsumers(bConsumers *nsqstore.BotConsumers) {
	bConsumers.Verification.AddHandler(bot)
}

func (bot *Bot) HandleNsqEvent(consumer *nsq.Consumer, cnf *config.NsqConfig) error {
	defer tools.Recovery(bot.logger)

	for {
		if err := consumer.ConnectToNSQLookupd(fmt.Sprintf("%s:%d", cnf.Host, cnf.Port)); err != nil {
			return err
		}
	}
}

func (bot *Bot) HandleBotEvent(ctx context.Context) error {
	defer tools.Recovery(bot.logger)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.botAPI.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		// Отсекать всех пользователей без username
		if !bot.check(bot.rewriter(update)) {
			continue
		}

		// Проверка на наличие незавершенных действий
		payload, ignore := bot.action(bot.rewriter(update))
		if !ignore && payload != nil {
			go bot.error(bot.rewriter(update), bot.ActionsHandler(ctx, update, payload))
			continue
		}

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Text {
			case static.BOT__CMD__START:
				go bot.error(update, bot.cmd.User().Start(ctx, update))
				continue

			case static.BOT__CMD__SKIP:
				go bot.error(update, bot.CancelAnyAction(ctx, update, payload))
				continue

			case static.BOT__CMD__DEV:
				go bot.error(update, bot.cmd.Base().Dev(ctx, update))
				continue

			case static.BOT__CMD__HELP:
				go bot.error(update, bot.cmd.Base().Help(ctx, update))
				continue

			default:
				continue
			}
		}

		if update.Message != nil {
			switch update.Message.Text {
			case static.BOT__BTN__BASE__NEW_EXCHANGE:
				go bot.error(update, bot.m.Exchange().NewExchange(ctx, update))
				continue

			case static.BOT__BTN__BASE__MY_BILLS:
				go bot.error(update, bot.m.Bill().MyBills(ctx, update))
				continue

			case static.BOT__BTN__OP__CANCEL:
				go bot.error(update, bot.CancelAnyAction(ctx, update, payload))
				continue

			case static.BOT__BTN__BASE__SUPPORT:
				go bot.error(update, bot.m.Base().SupportRequest(ctx, update))
				continue

			case static.BOT__BTN__BASE__ABOUT_BOT:
				go bot.error(update, bot.m.Base().AboutBot(ctx, update))
				continue

			case static.BOT__BTN__BASE__OPERATORS:
				go bot.error(update, bot.m.Base().Operators(ctx, update))
				continue

			default:
				continue
			}
		}

		if update.CallbackQuery != nil {
			fmt.Println(update.CallbackQuery.Data)
			// Декодирую полезную нагрузку
			p := map[string]interface{}{}
			bot.error(bot.rewriter(update), json.Unmarshal([]byte(update.CallbackQuery.Data), &p))

			switch p["CbQ"] {
			// Обработчики событий связанных с пользовательскими счетами
			case static.BOT__CQ_BL__ADD_BILL_S_1:
				go bot.error(bot.rewriter(update), bot.m.Bill().AddNewBillStepOne(ctx, update))
				continue

			// Обработчики событий связанных с операцией нового обмена
			case static.BOT__CQ__EX__COINS_TO_EXCHAGE:
				go bot.error(bot.rewriter(update), bot.m.Exchange().NewExchange(ctx, bot.rewriter(update)))
				continue

			case static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE:
				go bot.error(bot.rewriter(update), bot.m.Exchange().ReceiveAsResultOfExchange(ctx, update, p))
				continue

			case static.BOT__CQ__EX__REQ_BILL:
				go bot.error(bot.rewriter(update), bot.m.Exchange().СhooseBill(ctx, update, p))

			case static.BOT__CQ__EX__REQ_AMOUNT:
				go bot.error(bot.rewriter(update), bot.m.Exchange().ReqAmount(ctx, update, p))
				continue

			default:
				continue

			}
		}
	}
	return nil
}
