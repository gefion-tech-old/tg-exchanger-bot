package btns

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var UserKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(NEW_EXCHANGE__BTN),
		tgbotapi.NewKeyboardButton(MY_BILLS__BTN),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(MY_EXCHANGES__BTN),
		tgbotapi.NewKeyboardButton(SUPPORT__BTN),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ABOUT_BOT__BTN),
		tgbotapi.NewKeyboardButton(OPERATORS__BTN),
	),
)
