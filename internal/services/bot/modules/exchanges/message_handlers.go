package exchanges

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// @Button BOT__BTN__BASE__NEW_EXCHANGE
// @CallbackQuery BOT__CQ__EX__COINS_TO_EXCHAGE
// Отдает пользователю клав-список монет которые он может поменять
func (m *ModExchanges) NewExchange(ctx context.Context, update tgbotapi.Update) error {
	/*
		<ВЫПОЛНЕНИЕ ЗАПРОСА НА ПОЛУЧЕНИЕ ДОСТУПНЫХ ВАЛЮТ ДЛЯ ОБМЕНА>
	*/

	arrE := []*models.Exchanger{}

	// Тестовые вставки
	arrE = append(arrE, &models.Exchanger{ID: 1, Name: "Qiwi"})
	arrE = append(arrE, &models.Exchanger{ID: 2, Name: "Litecoin"})
	arrE = append(arrE, &models.Exchanger{ID: 3, Name: "Dogecoin"})
	arrE = append(arrE, &models.Exchanger{ID: 4, Name: "Bitcoin"})

	if update.CallbackQuery != nil {
		rMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		m.bAPI.Send(rMsg)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери валюту, которую нужно обменять 👇")
	msg.ReplyMarkup = m.kbd.Exchange().ExchangeCoinsList(arrE)
	m.bAPI.Send(msg)
	return nil
}
