package commands

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Commands struct {
	bAPI *tgbotapi.BotAPI
	sAPI api.ApiI
	kbd  keyboards.KeyboardsI

	userCommands UserCommandsI
	baseCommands BaseCommandsI

	logger *logrus.Logger
}

type CommandsI interface {
	Base() BaseCommandsI
	User() UserCommandsI
}

func InitCommands(bAPI *tgbotapi.BotAPI, kbd keyboards.KeyboardsI, sAPI api.ApiI, l *logrus.Logger) CommandsI {
	return &Commands{
		bAPI:   bAPI,
		sAPI:   sAPI,
		kbd:    kbd,
		logger: l,
	}
}

func (c *Commands) Base() BaseCommandsI {
	if c.baseCommands != nil {
		return c.baseCommands
	}

	c.baseCommands = InitBaseCommands(c.bAPI, c.sAPI, c.kbd, c.logger)
	return c.baseCommands
}

func (c *Commands) User() UserCommandsI {
	if c.userCommands != nil {
		return c.userCommands
	}

	c.userCommands = InitUserCommands(c.bAPI, c.sAPI, c.kbd, c.logger)
	return c.userCommands
}
