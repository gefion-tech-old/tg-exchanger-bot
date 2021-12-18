package commands

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type BaseCommands struct {
	bAPI   *tgbotapi.BotAPI
	sAPI   api.ApiI
	kbd    keyboards.KeyboardsI
	logger *logrus.Logger
}

type BaseCommandsI interface {
	Help(ctx context.Context, update tgbotapi.Update) error
	Dev(ctx context.Context, update tgbotapi.Update) error
}

func InitBaseCommands(bAPI *tgbotapi.BotAPI, sAPI api.ApiI, kbd keyboards.KeyboardsI, l *logrus.Logger) BaseCommandsI {
	return &BaseCommands{
		bAPI:   bAPI,
		sAPI:   sAPI,
		kbd:    kbd,
		logger: l,
	}
}

// @Command /help
func (c *BaseCommands) Help(ctx context.Context, update tgbotapi.Update) error {
	defer tools.Recovery(c.logger)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Какая-то информация тут.")
	msg.ParseMode = tgbotapi.ModeMarkdown
	c.bAPI.Send(msg)
	return nil
}

// @Command /dev
func (c *BaseCommands) Dev(ctx context.Context, update tgbotapi.Update) error {
	defer tools.Recovery(c.logger)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*Crafted by* [@I0HuKc](https://t.me/I0HuKc)")
	msg.ParseMode = tgbotapi.ModeMarkdown
	c.bAPI.Send(msg)
	return nil
}
