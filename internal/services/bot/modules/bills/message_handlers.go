package bills

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

func (m *ModBills) MyBills(ctx context.Context, update tgbotapi.Update) error {
	defer tools.Recovery(m.logger)

	// Вызываю через повторитель метод получения счетов пользователя
	r := api.Retry(m.sAPI.Bill().GetAll, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"chat_id": update.Message.From.ID,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		m.bAPI.Send(msg)
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case http.StatusOK:
		bills := []models.Bill{}

		if err := json.Unmarshal(resp.Body(), &bills); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список ваших счетов:")

		if len(bills) < 1 {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет добавленных счетов")
		}

		msg.ReplyMarkup = m.kbd.Bill().MyBillsList(bills)
		m.bAPI.Send(msg)
	}

	return nil
}
