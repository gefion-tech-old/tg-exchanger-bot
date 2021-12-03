package commands

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	START  = `/start`
	HELP   = `/help`
	AUTH   = `/auth`
	WHOAME = `/whoami`
)

type Commands struct {
	botAPI       *tgbotapi.BotAPI
	sAPI         api.ApiI
	userCommands *UserCommands
}

type CommandsI interface {
	User() UserCommandsI
}

func Init(bAPI *tgbotapi.BotAPI, sAPI api.ApiI) CommandsI {
	return &Commands{
		botAPI: bAPI,
		sAPI:   sAPI,
	}
}

func (c *Commands) User() UserCommandsI {
	if c.userCommands != nil {
		return c.userCommands
	}

	c.userCommands = &UserCommands{
		botAPI: c.botAPI,
		sAPI:   c.sAPI,
	}

	return c.userCommands
}
