package exchanges

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// @CallbackQuery BOT__CQ__EX__REQ_AMOUNT
func (m *ModExchanges) ReqAmount(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error {
	rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	m.bAPI.Send(rMsg)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Напиши сумму, которую хочешь обменять")
	msg.ReplyMarkup = m.kbd.Exchange().ReqAmountOffers()
	m.bAPI.Send(msg)
	return nil
}

// @CallbackQuery BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE
func (m *ModExchanges) ReceiveAsResultOfExchange(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error {
	/*
		<ВЫПОЛНЕНИЕ ЗАПРОСА НА ПОЛУЧЕНИЕ ДОСТУПНЫХ ВАЛЮТ ДЛЯ ОБМЕНА>
	*/

	arrE := []*models.Exchanger{}

	// Тестовые вставки
	arrE = append(arrE, &models.Exchanger{ID: 1, Name: "Сбербанк"})
	arrE = append(arrE, &models.Exchanger{ID: 2, Name: "Тинькофф"})
	arrE = append(arrE, &models.Exchanger{ID: 3, Name: "Альфа-Банк"})

	rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	m.bAPI.Send(rMsg)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Какую валюту хочешь получить?")
	msg.ReplyMarkup = m.kbd.Exchange().ReceiveAsResultOfExchangeList(arrE)
	m.bAPI.Send(msg)
	return nil
}
