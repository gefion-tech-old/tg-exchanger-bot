package commands

import (
	"context"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/btns"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/helpers"
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
	cHelloNewUserMsg := make(chan tgbotapi.MessageConfig)
	cHelloUserMsg := make(chan tgbotapi.MessageConfig)

	// Подгружаю сообщение для нового пользователя
	go func() {
		defer close(cHelloNewUserMsg)
		msg := helpers.GetMessage(ctx, update, c.sAPI, "hello_msg_new_user", update.Message.From.FirstName)
		cHelloNewUserMsg <- msg
	}()

	// Подгружаю сообщение для уже добавленого пользователя
	go func() {
		defer close(cHelloUserMsg)
		msg := helpers.GetMessage(ctx, update, c.sAPI, "hello_msg_user", update.Message.From.FirstName)
		cHelloUserMsg <- msg
	}()

	// Записываю в контекст UserReq
	ctx = context.WithValue(ctx, api.UserReqStructCtxKey, &models.UserReq{
		ChatID:   update.Message.From.ID,
		Username: update.Message.From.UserName,
	})

	// Вызываю через повторитель метод регистрации пользователя
	r := api.Retry(c.sAPI.User().Registration, 3, time.Second)
	resp, err := r(ctx)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		c.botAPI.Send(msg)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case http.StatusCreated:
		msg := <-cHelloNewUserMsg
		msg.ReplyMarkup = btns.UserKeyboard
		c.botAPI.Send(msg)
	case http.StatusUnprocessableEntity:
		msg := <-cHelloUserMsg
		msg.ReplyMarkup = btns.UserKeyboard
		c.botAPI.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Какая-то ошибка")
		c.botAPI.Send(msg)
	}
}
