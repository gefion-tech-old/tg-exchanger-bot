package keyboards

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ BaseKeyboardsI = (*BaseKeyboards)(nil)

type BaseKeyboards struct{}

type BaseKeyboardsI interface {
	CancelAction() tgbotapi.ReplyKeyboardMarkup
	BaseStartReplyMarkup() tgbotapi.ReplyKeyboardMarkup
}

// Отменить любое начатое действие
func (kb *BaseKeyboards) CancelAction() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(static.BOT__BTN__OP__CANCEL),
		),
	)
}

func (kb *BaseKeyboards) BaseStartReplyMarkup() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(static.BOT__BTN__BASE__NEW_EXCHANGE),
			tgbotapi.NewKeyboardButton(static.BOT__BTN__BASE__MY_BILLS),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(static.BOT__BTN__BASE__MY_EXCHANGES),
			tgbotapi.NewKeyboardButton(static.BOT__BTN__BASE__SUPPORT),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(static.BOT__BTN__BASE__ABOUT_BOT),
			tgbotapi.NewKeyboardButton(static.BOT__BTN__BASE__OPERATORS),
		),
	)
}
