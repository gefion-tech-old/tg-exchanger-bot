package keyboards

import (
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/core/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ ExchangeKeyboardsI = (*ExchangeKeyboards)(nil)

type ExchangeKeyboards struct{}

type ExchangeKeyboardsI interface {
	// InlineKeyboards
	ExchangeCoinsList(arr []*models.Coin) tgbotapi.InlineKeyboardMarkup
	ReceiveAsResultOfExchangeList(arr []*models.Coin, from string) tgbotapi.InlineKeyboardMarkup
	ReqAmountOffers(from string) tgbotapi.InlineKeyboardMarkup
	PayPage(url string) tgbotapi.InlineKeyboardMarkup
	СhooseBill(arr []models.Bill, from, to string) tgbotapi.InlineKeyboardMarkup
}

//  Список пользовательских счетов для выбора с какого проводить обмен
func (kb *ExchangeKeyboards) СhooseBill(arr []models.Bill, from, to string) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	for i := 0; i < len(arr); i++ {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Bill,
				fmt.Sprintf(`{"CbQ": "%s", "ID": %d, "From": "%s", "To": "%s"}`, static.BOT__CQ__EX__REQ_AMOUNT, arr[i].ID, from, to),
			),
		))
	}

	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf(`{"CbQ": "%s",  "From": "%s"}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE, from)),
		))

	return k
}

func (kb *ExchangeKeyboards) PayPage(url string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Перейти на страницу", url),
		),
	)
}

// Клавиатура для вывода списка валют, которых можно ПОМЕНЯТЬ
func (kb *ExchangeKeyboards) ExchangeCoinsList(arr []*models.Coin) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup

	for i := 0; i < len(arr); {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Name,
				fmt.Sprintf(`{"CbQ": "%s", "From": "%s"}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE, arr[i].ShortName),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i+1].Name,
				fmt.Sprintf(`{"CbQ": "%s", "From": "%s"}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE, arr[i+1].ShortName),
			),
		))
		i += 2
	}
	return k
}

// Клавиатура для вывода списка валют, которых можно ПОЛУЧИТЬ
func (kb *ExchangeKeyboards) ReceiveAsResultOfExchangeList(arr []*models.Coin, from string) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	for i := 0; i < len(arr); i++ {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Name,
				fmt.Sprintf(`{"CbQ": "%s", "From": "%s", "To": "%s", "F": %t}`, static.BOT__CQ__EX__REQ_BILL, from, arr[i].ShortName, arr[i].Fiat),
			),
		))
	}

	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ__EX__COINS_TO_EXCHAGE)),
		))

	return k
}

func (kb *ExchangeKeyboards) ReqAmountOffers(from string) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf(`{"CbQ": "%s", "From": "%s"}`, static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE, from)),
		))

	return k
}
