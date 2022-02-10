package exchanges

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/core/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

func (m *ModExchanges) ReqAmountForCrypto(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {

	// Получение обменника
	r := api.Retry(m.sAPI.Exchanger().Get, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"name": "1obmen",
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		m.bAPI.Send(msg)
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != http.StatusOK {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Не удалось получить данные актуального курса ❌")
		m.bAPI.Send(msg)
		return nil
	}

	body := models.Exchanger{}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return err
	}

	return nil
}

// Функция вызывается если клиент хочет получить деньги в крипте
// следовательно необходимо получить от пользователя адрес кошелька
// на который он желает получить конвертируемые средства
func (m *ModExchanges) HandleReceivedAddress(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
	defer tools.Recovery(m.logger)

	if update.Message.Text != "" {
		q, err := m.GetExchangeInfo(ctx, update, action.MetaData["From"].(string), action.MetaData["To"].(string))
		if err != nil {
			return err
		}

		// Обновление в redis пользовательского действия
		if err := m.redis.UserActions().New(update.Message.Chat.ID, &models.UserAction{
			ActionType: static.BOT__A__EX__NEW_EXCHAGE,
			Step:       1,
			MetaData: map[string]interface{}{
				"From":      action.MetaData["From"],
				"To":        action.MetaData["To"],
				"ToFiat":    action.MetaData["ToFiat"],
				"Bill":      update.Message.Text,
				"Course":    q.In,
				"MinAmount": strings.Split(q.MinAmount, " ")[0],
				"MaxAmount": strings.Split(q.MaxAmount, " ")[0],
			},
			User: struct {
				ChatID   int
				Username string
			}{
				ChatID:   int(update.Message.Chat.ID),
				Username: update.Message.Chat.UserName,
			},
		}); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🟡Получение актуального курса🟡")
		waitM, _ := m.bAPI.Send(msg)

		rMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, waitM.MessageID)
		m.bAPI.Send(rMsg)

		text := fmt.Sprintf("Напиши сумму обмена 👇\n\n*От*: `%s`\n*До*: `%s`\n*Курс*: `%0.3f`", q.MinAmount, q.MaxAmount, q.In)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown
		m.bAPI.Send(msg)
		return nil

	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Хм, это непохоже на адрес кошелька...")
	m.bAPI.Send(msg)
	return nil
}

func (m *ModExchanges) CreateLinkForPayment(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
	defer tools.Recovery(m.logger)

	if update.Message.Text != "" {
		// Регулярка для получения дробного значения из текста
		re := regexp.MustCompile(`(?:\d+(?:\.\d*)?|\.\d+)`)

		max, err := strconv.ParseFloat(re.FindAllString(action.MetaData["MaxAmount"].(string), -1)[0], 64)
		if err != nil {
			return err
		}

		min, err := strconv.ParseFloat(re.FindAllString(action.MetaData["MinAmount"].(string), -1)[0], 64)
		if err != nil {
			return err
		}

		need, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Используйте числовые значения ❌")
			m.bAPI.Send(msg)
			return nil
		}

		// Определяю допустимо ли размер запрашиваемого обмена
		if need > max || need < min {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Недопустимая сумма обмена ❌")
			m.bAPI.Send(msg)
			return nil
		}

		r := api.Retry(m.sAPI.Exchanger().GetAdress, 3, time.Second)
		resp, err := r(ctx, map[string]interface{}{
			"merchant": "whitebit",

			"exchange_from":   action.MetaData["From"],
			"exchange_to":     action.MetaData["To"],
			"course":          fmt.Sprintf("%f", action.MetaData["Course"]),
			"expected_amount": 10,
			"client_address":  "3MGgZg2k1bKd1n598xewrDsCdYUfi3JWgu",
			"created_by": map[string]interface{}{
				"username": update.Message.Chat.UserName,
				"chat_id":  update.Message.Chat.ID,
			},
		})
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
			m.bAPI.Send(msg)
			return err
		}
		defer fasthttp.ReleaseResponse(resp)

		if resp.StatusCode() != http.StatusOK {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Не удалось получить адрес ❌")
			m.bAPI.Send(msg)
			return nil
		}

		body := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body(), &body); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Хорошо, адрес для перевода 👇\n\n`%s`", body["account"].(map[string]interface{})["address"]))
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = m.kbd.Base().BaseStartReplyMarkup()
		m.bAPI.Send(msg)

		return m.redis.UserActions().Delete(update.Message.Chat.ID)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хм, это непохоже на сумму...")
	m.bAPI.Send(msg)
	return nil
}
