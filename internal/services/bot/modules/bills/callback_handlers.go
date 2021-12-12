package bills

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	validation "github.com/go-ozzo/ozzo-validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

func (m *ModBills) AddNewBillStepFour(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
	if update.Message.Photo != nil {
		// Вызываю через повторитель метод отправки уведомления на сервер
		r := api.Retry(m.sAPI.Notification().Create, 3, time.Second)
		resp, err := r(ctx, map[string]interface{}{
			"type": action.ActionType,
			"meta_data": map[string]interface{}{
				"code":      action.MetaData["Code"],
				"user_card": action.MetaData["Card"],
			},
			"user": map[string]interface{}{
				"chat_id":  action.User.ChatID,
				"username": action.User.Username,
			},
		})
		if err != nil {
			return errors.ErrBotServerNoAnswer
		}
		defer fasthttp.ReleaseResponse(resp)

		switch resp.StatusCode() {
		case http.StatusCreated:
			if err := m.redis.UserActions().Delete(update.Message.Chat.ID); err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Заявка принята ✅\n\nПосле верификации карты менеджером, вам будет отправленно уведомление.")
			msg.ParseMode = tgbotapi.ModeMarkdown
			msg.ReplyMarkup = m.kbd.Base().BaseStartReplyMarkup()
			m.bAPI.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Похоже произошла какая-та ошибка, попробуйте повторить попытку.")
			m.bAPI.Send(msg)
			return nil
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хм, это не похоже на фотографию...")
	m.bAPI.Send(msg)
	return nil
}

// Если пользователь ввел не ту карту и хочет отменить этот шаг
func (m *ModBills) AddNewBillStepThreeInCorrect(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
	action.MetaData = nil
	action.Step--
	if err := m.redis.UserActions().New(update.CallbackQuery.Message.Chat.ID, action); err != nil {
		return err
	}

	return m.AddNewBillStepOne(ctx, update)
}

func (m *ModBills) AddNewBillStepThreeCorrect(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
	c := tools.VerificationCode(false)
	action.MetaData["Code"] = c
	action.Step++

	if err := m.redis.UserActions().New(update.CallbackQuery.Message.Chat.ID, action); err != nil {
		return err
	}

	rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	m.bAPI.Send(rMsg)

	text := fmt.Sprintf("Хорошо, пришли мне фото с кодом подтверждения.\nТут текст о том как нужно сфотографировать карту с кодом подтверждения.\n\nКод подтверждения: *%d*", c)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	m.bAPI.Send(msg)
	return nil
}

func (m *ModBills) AddNewBillStepTwo(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
	// Валидация номера карты
	pattern := `^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`
	if err := validation.Validate(update.Message.Text, validation.Required, validation.Match(regexp.MustCompile(pattern))); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Неверный формат ❌")
		m.bAPI.Send(msg)
		return nil
	}

	// Обновляю запись о пользовательском действии
	action.MetaData = map[string]interface{}{"Card": update.Message.Text}
	action.Step++
	if err := m.redis.UserActions().New(update.Message.Chat.ID, action); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Добавляем карту `%s`.\nВсе верно?", update.Message.Text))
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = m.kbd.Bill().CardСorrectnessConfirmation()
	m.bAPI.Send(msg)
	return nil
}

// Запрос у пользователя номера карты
func (m *ModBills) AddNewBillStepOne(ctx context.Context, update tgbotapi.Update) error {

	if err := m.redis.UserActions().New(update.CallbackQuery.Message.Chat.ID, &models.UserAction{
		ActionType: static.BOT__A__BL__ADD_NEW_BILL,
		Step:       1,
		User: models.UserReq{
			ChatID:   int(update.CallbackQuery.Message.Chat.ID),
			Username: update.CallbackQuery.Message.Chat.UserName,
		},
	}); err != nil {
		return err
	}

	rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	m.bAPI.Send(rMsg)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Хорошо, пришли мне номер карты которую необходимо добавить 👇")
	msg.ReplyMarkup = m.kbd.Base().CancelAction()
	m.bAPI.Send(msg)
	return nil
}