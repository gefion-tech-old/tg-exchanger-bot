package bot

import (
	"context"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/commands"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
	cmd    commands.CommandsI
}

type BotI interface {
	/*
		Метод слушатель входящих сообщений из очереди
	*/
	HandleNsqEvent(consumer *nsq.Consumer, cnf *config.NsqConfig) error
	/*
		Метод слушатель входящих событий в telegram
	*/
	HandleBotEvent(ctx context.Context) error
	/*
		Коннектор всех nsq потребителей
	*/
	ConnectNsqConsumers(bConsumers *nsqstore.BotConsumers)
}

func Init(bAPI *tgbotapi.BotAPI, sAPI api.ApiI) BotI {
	return &Bot{
		botAPI: bAPI,
		cmd:    commands.Init(bAPI, sAPI),
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
		if update.Message != nil {
			// Отсекать всех пользователей без username
			if update.Message.From.UserName == "" {
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Привет! Вам необходимо установить себе `username` для использования этого бота. После этого выполните команду /start чтобы перезапустить бота.")
				msg.ParseMode = tgbotapi.ModeMarkdown
				bot.botAPI.Send(msg)
				continue
			}

			switch update.Message.Text {
			case commands.START__CMD:
				go func() {
					bot.cmd.User().Start(ctx, update)
				}()
			default:
				continue
			}
		}
	}
	return nil
}
