package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/gefion-tech/tg-exchanger-bot/internal/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run bot",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Init()
			if _, err := toml.DecodeFile(
				fmt.Sprintf("config/config.%s.toml", viper.GetString("env")), cfg); err != nil {
				panic(err)
			}

			if err := runner(cfg); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().IntP("cpu", "c", 2, "Number of processor threads")
	cmd.Flags().StringP("env", "e", "local", "Launch environment")

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}

func runner(cfg *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Инициализация приложения

	logger := logrus.New()

	errs, ctx := errgroup.WithContext(ctx)

	// Инициализация redis хранилищ

	uActionsClient, err := db.InitRedis(&cfg.Redis, 10)
	if err != nil {
		return err
	}

	uActions := redisstore.InitUserActions(uActionsClient)
	defer uActions.Close()

	// Инициализация общей сборки всех redis хранилищ используемых в этом приложении
	aRedis := redisstore.InitRedisStore(uActions)

	botAPI, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return err
	}

	botAPI.Debug = cfg.Bot.Debug

	// Инициализация модуля работы с API сервера
	sAPI := api.Init(&cfg.API, &cfg.Bot, logger)

	// Инициализирую модуль бота
	bot := bot.Init(botAPI, sAPI, aRedis, &cfg.Bot, logger)

	// Инициализирую всех NSQ потребителей
	bConsumers, teardown, err := nsqstore.Init(&cfg.NSQ)
	if err != nil {
		return err
	}
	defer teardown(bConsumers.Verification)

	// Подключаю всех NSQ потребителей
	bot.ConnectNsqConsumers(bConsumers)

	// Запуск обработчиков всех событий
	errs.Go(func() error { return bot.HandleNsqEvent(bConsumers.Verification, &cfg.NSQ) })
	errs.Go(func() error { return bot.HandleBotEvent(ctx) })

	<-ctx.Done()
	stop()

	if errs.Wait() != nil {
		fmt.Println(errs.Wait().Error())
	}

	return nil
}
