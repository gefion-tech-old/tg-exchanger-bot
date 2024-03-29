package commands

import (
	"context"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/helpers"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type UserCommands struct {
	bAPI   *tgbotapi.BotAPI
	sAPI   api.ApiI
	kbd    keyboards.KeyboardsI
	logger *logrus.Logger
}

type UserCommandsI interface {
	Start(ctx context.Context, update tgbotapi.Update) error
}

func InitUserCommands(bAPI *tgbotapi.BotAPI, sAPI api.ApiI, kbd keyboards.KeyboardsI, l *logrus.Logger) UserCommandsI {
	return &UserCommands{
		bAPI:   bAPI,
		sAPI:   sAPI,
		kbd:    kbd,
		logger: l,
	}
}

// @Command /start
func (c *UserCommands) Start(ctx context.Context, update tgbotapi.Update) error {
	defer tools.Recovery(c.logger)

	errs, _ := errgroup.WithContext(ctx)

	cHelloNewUserMsg := make(chan *models.Message)
	cHelloUserMsg := make(chan *models.Message)

	// Подгружаю сообщение для нового пользователя
	errs.Go(func() error {
		defer close(cHelloNewUserMsg)
		msg, err := helpers.GetMessage(ctx, update, c.sAPI, "hello_msg_new_user", c.logger, update.Message.From.FirstName)
		if err != nil {
			return err
		}

		cHelloNewUserMsg <- msg
		return nil
	})

	// Подгружаю сообщение для уже добавленого пользователя
	errs.Go(func() error {
		defer close(cHelloUserMsg)
		msg, err := helpers.GetMessage(ctx, update, c.sAPI, "hello_msg_user", c.logger, update.Message.From.FirstName)
		if err != nil {
			return err
		}

		cHelloUserMsg <- msg
		return nil
	})

	// Вызываю через повторитель метод регистрации пользователя
	r := api.Retry(c.sAPI.User().Registration, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"chat_id":  update.Message.From.ID,
		"username": update.Message.From.UserName,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		c.bAPI.Send(msg)
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	// Дожидаюсь ответа всех горутин
	helloNewUserMsg := <-cHelloNewUserMsg
	helloUserMsg := <-cHelloUserMsg

	// Все горутины завершились успешно
	if helloNewUserMsg == nil || helloUserMsg == nil {
		return errs.Wait()
	}

	switch resp.StatusCode() {
	case http.StatusCreated:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloNewUserMsg.MessageText)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = c.kbd.Base().BaseStartReplyMarkup()
		c.bAPI.Send(msg)
		return nil
	case http.StatusUnprocessableEntity:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloUserMsg.MessageText)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = c.kbd.Base().BaseStartReplyMarkup()
		c.bAPI.Send(msg)
		return nil
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Какая-то ошибка")
		c.bAPI.Send(msg)
		return nil
	}
}
