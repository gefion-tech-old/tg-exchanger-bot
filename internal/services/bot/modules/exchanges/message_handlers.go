package exchanges

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// @Button BOT__BTN__BASE__NEW_EXCHANGE
// @CallbackQuery BOT__CQ__EX__COINS_TO_EXCHAGE
// Отдает пользователю клав-список монет которые он может поменять
func (m *ModExchanges) NewExchange(ctx context.Context, update tgbotapi.Update) error {
	defer tools.Recovery(m.logger)

	/*
		<ВЫПОЛНЕНИЕ ЗАПРОСА НА ПОЛУЧЕНИЕ ДОСТУПНЫХ ВАЛЮТ ДЛЯ ОБМЕНА>
	*/

	if update.CallbackQuery != nil {
		rMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		m.bAPI.Send(rMsg)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери валюту, которую нужно обменять 👇")
	msg.ReplyMarkup = m.kbd.Exchange().ExchangeCoinsList(models.COINS)
	m.bAPI.Send(msg)
	return nil
}
