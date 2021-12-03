package app

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot"
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

	// Инициализирую модуль бота
	bot := bot.Init(botAPI, sAPI)
	return bot.EventHandler(ctx)
}
