package bot

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
	cmd    commands.CommandsI
}

type BotI interface {
	/*
		Метод слушатель входящих событий
	*/
	EventHandler(ctx context.Context) error
}

func Init(bAPI *tgbotapi.BotAPI, sAPI api.ApiI) BotI {
	return &Bot{
		botAPI: bAPI,
		cmd:    commands.Init(bAPI, sAPI),
	}
}

func (bot *Bot) EventHandler(ctx context.Context) error {
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
			case commands.START:
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
