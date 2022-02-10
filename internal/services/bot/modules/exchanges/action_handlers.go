package exchanges

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

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

			"exchange_from":   "USDTTRC20",
			"exchange_to":     "SBERRUB",
			"course":          "76.0947",
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
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Не удалось получить данные актуального курса ❌")
			m.bAPI.Send(msg)
			return nil
		}

		body := map[string]interface{}{}
		if err := json.Unmarshal(resp.Body(), &body); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Хорошо, адрес для перевода 👇\n\n`%s`", body["account"].(map[string]interface{})["address"]))
		msg.ParseMode = tgbotapi.ModeMarkdown
		// msg.ReplyMarkup = m.kbd.Exchange().PayPage("https://some.com")
		m.bAPI.Send(msg)
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хм, это непохоже на сумму...")
	m.bAPI.Send(msg)
	return nil
}
