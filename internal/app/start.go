package app

import (
	"context"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type App struct {
	config *config.Config
}

type AppI interface {
	Start(ctx context.Context) error
}

func Init(c *config.Config) AppI {
	return &App{
		config: c,
	}
}

func (a *App) Start(ctx context.Context) error {
	botAPI, err := tgbotapi.NewBotAPI(a.config.Bot.Token)
	if err != nil {
		return err
	}

	botAPI.Debug = a.config.Bot.Debug

	// Инициализация модуля работы с API сервера
	sAPI := api.Init(&a.config.API)

	mEventConsumer, err := nsqstore.Init(&a.config.NSQ, "verification-code", "telegram")
	if err != nil {
		return err
	}
	mEventConsumer.AddHandler(&nsqstore.EventsHandler{BotAPI: botAPI})
	defer mEventConsumer.Stop()

	go func() {
		for {
			mEventConsumer.ConnectToNSQLookupd(fmt.Sprintf("%s:%d", a.config.NSQ.Host, a.config.NSQ.Port))
		}
	}()

	// Инициализирую модуль бота
	bot := bot.Init(botAPI, sAPI)
	return bot.MessageEventHandler(ctx)

}
