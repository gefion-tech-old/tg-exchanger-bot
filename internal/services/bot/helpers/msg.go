package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

func GetMessage(ctx context.Context, update tgbotapi.Update, sAPI api.ApiI, connector string, params ...interface{}) tgbotapi.MessageConfig {
	ctx = context.WithValue(ctx, api.MessageConnectorCtxKey, connector)
	r := api.Retry(sAPI.Message().Get, 3, time.Second)

	resp, err := r(ctx)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
	}
	defer fasthttp.ReleaseResponse(resp)

	var body map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка декодирования ответа...")
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(body["message_text"].(string), params))
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	case http.StatusNotFound:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сообщение не найдено...")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла какая-то ошибка, попробуйте повторить попытку позже.")
		msg.ParseMode = tgbotapi.ModeMarkdown
		return msg
	}
}
