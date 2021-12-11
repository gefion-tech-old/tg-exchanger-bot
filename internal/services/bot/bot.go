package bot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/commands"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
	cnf    *config.BotConfig
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

func Init(bAPI *tgbotapi.BotAPI, sAPI api.ApiI, cnf *config.BotConfig) BotI {
	kb := keyboards.InitKeyboards()
	mod := modules.InitBotModules(bAPI, kb, sAPI)
	cmd := commands.InitCommands(bAPI, kb, sAPI)

	return &Bot{
		botAPI: bAPI,
		cnf:    cnf,
		cmd:    cmd,
		kbd:    kb,
		m:      mod,
	}
}

func (bot *Bot) ConnectNsqConsumers(bConsumers *nsqstore.BotConsumers) {
	bConsumers.Verification.AddHandler(bot)
}

func (bot *Bot) HandleNsqEvent(consumer *nsq.Consumer, cnf *config.NsqConfig) error {
	for {
		if err := consumer.ConnectToNSQLookupd(fmt.Sprintf("%s:%d", cnf.Host, cnf.Port)); err != nil {
			return err
		}
	}
}

func (bot *Bot) HandleBotEvent(ctx context.Context) error {
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

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Text {
			case static.BOT__CMD__START:
				go bot.error(update, bot.cmd.User().Start(ctx, update))
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
			default:
				continue
			}
		}

		if update.CallbackQuery != nil {
			// Декодирую полезную нагрузку
			p := map[string]interface{}{}
			bot.error(update, json.Unmarshal([]byte(update.CallbackQuery.Data), &p))

			switch p["CbQ"] {
			case static.BOT__CQ__EX__COINS_TO_EXCHAGE:
				go bot.error(bot.rewriter(update), bot.m.Exchange().NewExchange(ctx, bot.rewriter(update)))
				continue
			case static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE:
				go bot.error(bot.rewriter(update), bot.m.Exchange().ReceiveAsResultOfExchange(ctx, update, p))
				continue
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
