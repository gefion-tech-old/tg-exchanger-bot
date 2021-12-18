package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	logger *logrus.Logger
	config *config.Config
}

type AppI interface {
	Start(ctx context.Context) error
}

func Init(c *config.Config) AppI {
	return &App{
		config: c,
		logger: logrus.New(),
	}
}

func (a *App) Start(ctx context.Context) error {
	f, err := os.Create(fmt.Sprintf("logs/%s.log", time.Now().UTC().Format("2006-01-02T15:04:05.00000000")))
	if err != nil {
		return err
	}

	a.logger.SetOutput(f)
	a.logger.SetLevel(logrus.ErrorLevel)

	errs, ctx := errgroup.WithContext(ctx)

	// Инициализация redis хранилищ

	uActionsClient, err := db.InitRedis(&a.config.Redis, 10)
	if err != nil {
		return err
	}

	uActions := redisstore.InitUserActions(uActionsClient)
	defer uActions.Close()

	// Инициализация общей сборки всех redis хранилищ используемых в этом приложении
	aRedis := redisstore.InitRedisStore(uActions)

	botAPI, err := tgbotapi.NewBotAPI(a.config.Bot.Token)
	if err != nil {
		return err
	}

	botAPI.Debug = a.config.Bot.Debug

	// Инициализация модуля работы с API сервера
	sAPI := api.Init(&a.config.API, &a.config.Bot, a.logger)

	// Инициализирую модуль бота
	bot := bot.Init(botAPI, sAPI, aRedis, &a.config.Bot, a.logger)

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
