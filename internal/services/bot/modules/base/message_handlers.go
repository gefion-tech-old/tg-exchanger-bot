package base

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

func (m *ModBase) SupportRequest(ctx context.Context, update tgbotapi.Update) error {
	// Вызываю через повторитель метод отправки уведомления на сервер
	r := api.Retry(m.sAPI.Notification().Create, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"type": static.BOT__A__BASE__REQ_SUPPORT,
		"user": map[string]interface{}{
			"chat_id":  update.Message.From.ID,
			"username": update.Message.From.UserName,
		},
	})
	if err != nil {
		return errors.ErrBotServerNoAnswer
	}
	defer fasthttp.ReleaseResponse(resp)

	fmt.Println(resp.StatusCode())

	switch resp.StatusCode() {
	case http.StatusCreated:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Заявка принята ✅\n\nНаши менеджеры скоро свяжутся с вами")
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = m.kbd.Base().BaseStartReplyMarkup()
		m.bAPI.Send(msg)
		return nil

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже произошла какая-та ошибка, попробуйте повторить попытку.")
		m.bAPI.Send(msg)
		return nil
	}
}
