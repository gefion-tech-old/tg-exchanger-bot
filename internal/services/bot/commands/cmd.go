package commands

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Commands struct {
	bAPI *tgbotapi.BotAPI
	sAPI api.ApiI
	kbd  keyboards.KeyboardsI

	userCommands *UserCommands
}

type CommandsI interface {
	User() UserCommandsI
}

func InitCommands(bAPI *tgbotapi.BotAPI, kbd keyboards.KeyboardsI, sAPI api.ApiI) CommandsI {
	return &Commands{
		bAPI: bAPI,
		sAPI: sAPI,
		kbd:  kbd,
	}
}

func (c *Commands) User() UserCommandsI {
	if c.userCommands != nil {
		return c.userCommands
	}

	c.userCommands = &UserCommands{
		bAPI: c.bAPI,
		sAPI: c.sAPI,
		kbd:  c.kbd,
	}

	return c.userCommands
}
