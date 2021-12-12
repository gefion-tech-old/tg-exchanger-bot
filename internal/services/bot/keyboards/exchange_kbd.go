package keyboards

import (
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ ExchangeKeyboardsI = (*ExchangeKeyboards)(nil)

type ExchangeKeyboards struct{}

type ExchangeKeyboardsI interface {
	// InlineKeyboards
	ExchangeCoinsList(arrE []*models.Exchanger) tgbotapi.InlineKeyboardMarkup
	ReceiveAsResultOfExchangeList(arr []*models.Exchanger) tgbotapi.InlineKeyboardMarkup
	ReqAmountOffers() tgbotapi.InlineKeyboardMarkup
}

// Клавиатура для вывода списка валют, которых можно ПОМЕНЯТЬ
func (kb *ExchangeKeyboards) ExchangeCoinsList(arr []*models.Exchanger) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup

	for i := 0; i < len(arr); {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Name,
				fmt.Sprintf(`{"CbQ": "%s", "ID": %d}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE, arr[i].ID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i+1].Name,
				fmt.Sprintf(`{"CbQ": "%s", "ID": %d}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE, arr[i+1].ID),
			),
		))
		i += 2
	}
	return k
}

// Клавиатура для вывода списка валют, которых можно ПОЛУЧИТЬ
func (kb *ExchangeKeyboards) ReceiveAsResultOfExchangeList(arr []*models.Exchanger) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	for i := 0; i < len(arr); i++ {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Name,
				fmt.Sprintf(`{"CbQ": "%s", "ID": %d}`, static.BOT__CQ__EX__REQ_AMOUNT, arr[i].ID),
			),
		))
	}

	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ__EX__COINS_TO_EXCHAGE)),
		))

	return k
}

func (kb *ExchangeKeyboards) ReqAmountOffers() tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE)),
		))

	return k
}