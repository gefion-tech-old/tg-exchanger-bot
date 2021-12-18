package exchanges

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
)

// @CallbackQuery BOT__CQ__EX__REQ_AMOUNT
func (m *ModExchanges) ReqAmount(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error {
	defer tools.Recovery(m.logger)

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

	body := map[string]interface{}{}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return err
	}

	// Создание в redis пользовательского действия
	if err := m.redis.UserActions().New(update.CallbackQuery.Message.Chat.ID, &models.UserAction{
		ActionType: static.BOT__A__EX__NEW_EXCHAGE,
		Step:       1,
		MetaData: map[string]interface{}{
			"From": p["From"],
			"To":   p["To"],
		},
		User: struct {
			ChatID   int
			Username string
		}{
			ChatID:   int(update.CallbackQuery.Message.Chat.ID),
			Username: update.CallbackQuery.Message.Chat.UserName,
		},
	}); err != nil {
		return err
	}

	rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID-1)
	m.bAPI.Send(rMsg)

	rMsg = tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	m.bAPI.Send(rMsg)

	msgInfo := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Обмен из *%s* в *%s*", p["From"], p["To"]))
	msgInfo.ParseMode = tgbotapi.ModeMarkdown
	msgInfo.ReplyMarkup = m.kbd.Base().CancelAction()
	m.bAPI.Send(msgInfo)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "🟡Получение актуального курса🟡")
	waitM, _ := m.bAPI.Send(msg)

	// Получение актуальных котировок
	data, err := m.quotes(ctx, update, body["url"].(string))
	if err != nil {
		return err
	}

	// Поиск нужной котировки
	cQ := make(chan *models.OneObmenItem)
	for i := 0; i < len(data.Rates); i++ {
		go func(i int) {
			if data.Rates[i].From == p["From"] && data.Rates[i].To == p["To"] {
				defer close(cQ)
				cQ <- &data.Rates[i]
			}
		}(i)
	}

	if q := <-cQ; q != nil {
		// Обновление в redis пользовательского действия
		if err := m.redis.UserActions().New(update.CallbackQuery.Message.Chat.ID, &models.UserAction{
			ActionType: static.BOT__A__EX__NEW_EXCHAGE,
			Step:       1,
			MetaData: map[string]interface{}{
				"From":      p["From"],
				"To":        p["To"],
				"MinAmount": q.MinAmount,
				"MaxAmount": q.MaxAmount,
			},
			User: struct {
				ChatID   int
				Username string
			}{
				ChatID:   int(update.CallbackQuery.Message.Chat.ID),
				Username: update.CallbackQuery.Message.Chat.UserName,
			},
		}); err != nil {
			return err
		}

		rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, waitM.MessageID)
		m.bAPI.Send(rMsg)

		text := fmt.Sprintf("Напиши сумму обмена 👇\n\n*От*: `%s`\n*До*: `%s`\n*Курс*: `%0.3f`", q.MinAmount, q.MaxAmount, q.In)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = m.kbd.Exchange().ReqAmountOffers(p["From"].(string))
		m.bAPI.Send(msg)
		return nil
	}

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Не удалось найти курс для данного направления обмена. Попробуйте повторить попытку позже или выбрать другое направление обмена.")
	m.bAPI.Send(msg)
	return nil
}

// @CallbackQuery BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE
func (m *ModExchanges) ReceiveAsResultOfExchange(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error {
	defer tools.Recovery(m.logger)

	if err := m.redis.UserActions().Delete(update.CallbackQuery.Message.Chat.ID); err != nil {
		return err
	}

	// Определение какие поддерживаются направления обмена для выбранной валюты
	directions := []*models.Direction{}
	for i := 0; i < len(models.DIRECTIONS); i++ {
		if models.DIRECTIONS[i].From == p["From"].(string) {
			directions = append(directions, models.DIRECTIONS[i])
		}
	}

	coins := []*models.Coin{}
	for d := 0; d < len(directions); d++ {
		for c := 0; c < len(models.COINS); c++ {
			if directions[d].To == models.COINS[c].ShortName {
				coins = append(coins, models.COINS[c])
			}
		}
	}

	rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	m.bAPI.Send(rMsg)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Обмен из *%s*", p["From"]))
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = m.kbd.Base().BaseStartReplyMarkup()
	m.bAPI.Send(msg)

	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "На данный момент по этой валюте нет поддерживаемых направлений обмена.")

	if len(coins) > 0 {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Какую валюту хочешь получить?")
		msg.ReplyMarkup = m.kbd.Exchange().ReceiveAsResultOfExchangeList(coins, p["From"].(string))
	}

	m.bAPI.Send(msg)
	return nil
}

// Вспомогательный метод
// Запрашивает актуальные котировки валют
// Находит нужную котировку и возвращает ее
func (m *ModExchanges) quotes(ctx context.Context, update tgbotapi.Update, url string) (*models.OneObmen, error) {
	defer tools.Recovery(m.logger)

	r := api.Retry(m.sAPI.Exchanger().GetQuotesXML, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"url": url,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		m.bAPI.Send(msg)
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	data := models.OneObmen{}
	if err := xml.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}

	return &data, nil
}
