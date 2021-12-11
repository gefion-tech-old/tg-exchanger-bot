package app

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/sync/errgroup"
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
	errs, ctx := errgroup.WithContext(ctx)

	botAPI, err := tgbotapi.NewBotAPI(a.config.Bot.Token)
	if err != nil {
		return err
	}

	botAPI.Debug = a.config.Bot.Debug

	// Инициализация модуля работы с API сервера
	sAPI := api.Init(&a.config.API)

	// Инициализирую модуль бота
	bot := bot.Init(botAPI, sAPI, &a.config.Bot)

	// Инициализирую всех NSQ потребителей
	bConsumers, teardown, err := nsqstore.Init(&a.config.NSQ)
	if err != nil {
		return err
	}
	defer teardown(bConsumers.Verification)

	// Подключаю всех NSQ потребителей
	bot.ConnectNsqConsumers(bConsumers)

	// Запуск обработчиков всех событий
	errs.Go(func() error { return bot.HandleNsqEvent(bConsumers.Verification, &a.config.NSQ) })
	errs.Go(func() error { return bot.HandleBotEvent(ctx) })

	return errs.Wait()
}
