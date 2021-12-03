package commands

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/btns"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

type UserCommands struct {
	botAPI *tgbotapi.BotAPI
	sAPI   api.ApiI
}

type UserCommandsI interface {
	Start(ctx context.Context, update tgbotapi.Update)
}

func (c *UserCommands) Start(ctx context.Context, update tgbotapi.Update) {
	msg := tgbotapi.MessageConfig{}
	ur := &models.UserReq{
		ChatID:   update.Message.From.ID,
		Username: update.Message.From.UserName,
	}

	// Записываю в контекст UserReq
	ctx = context.WithValue(ctx, api.UserReqStructCtxKey, ur)

	// Вызываю через повторитель метод регистрации пользователя
	r := api.Retry(c.sAPI.User().Registration, 3, time.Second)
	resp, err := r(ctx)
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		c.botAPI.Send(msg)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case http.StatusCreated:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s!\nВы успешно зарегестрировались в боте.", update.Message.From.FirstName))
		msg.ReplyMarkup = btns.UserKeyboard
	case http.StatusUnprocessableEntity:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("С возвращением, %s!", update.Message.From.FirstName))
		msg.ReplyMarkup = btns.UserKeyboard
	default:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Какая-то ошибка")
	}

	c.botAPI.Send(msg)
}
