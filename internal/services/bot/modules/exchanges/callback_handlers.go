package exchanges

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/core/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

func (m *ModExchanges) СhooseBill(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error {
	defer tools.Recovery(m.logger)

	// Вызываю через повторитель метод получения счетов пользователя
	r := api.Retry(m.sAPI.Bill().GetAll, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"chat_id": update.CallbackQuery.Message.Chat.ID,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Сервер не отвечает")
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

		rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
		m.bAPI.Send(rMsg)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите счет:")

		if len(bills) == 0 {
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите счет: \n\nУ вас нет добавленных счетов")
		}

		msg.ReplyMarkup = m.kbd.Exchange().СhooseBill(bills, p["From"].(string), p["To"].(string))
		m.bAPI.Send(msg)
	}

	return nil
}

// @CallbackQuery BOT__CQ__EX__REQ_AMOUNT
func (m *ModExchanges) ReqAmount(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error {
	defer tools.Recovery(m.logger)

	errs, _ := errgroup.WithContext(ctx)

	cE := make(chan *models.Exchanger, 1)
	cB := make(chan *models.Bill, 1)

	// Получение информации по выбранному счету
	errs.Go(func() error {
		defer close(cB)

		// Вызываю через повторитель метод получения счетов пользователя
		r := api.Retry(m.sAPI.Bill().GetBill, 3, time.Second)
		resp, err := r(ctx, map[string]interface{}{
			"bill_id": int(p["ID"].(float64)),
		})
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Сервер не отвечает")
			m.bAPI.Send(msg)
			return err
		}
		defer fasthttp.ReleaseResponse(resp)

		if resp.StatusCode() != http.StatusOK {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Не удалось получить по выбранному счету ❌")
			m.bAPI.Send(msg)
			return nil
		}

		body := models.Bill{}
		if err := json.Unmarshal(resp.Body(), &body); err != nil {
			return err
		}

		cB <- &body
		return nil
	})

	// Получение информации по  обменнику
	errs.Go(func() error {
		defer close(cE)

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

		cE <- &body
		return nil
	})

	e := <-cE
	b := <-cB

	if e == nil || b == nil {
		return errs.Wait()
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
	data, err := m.quotes(ctx, update, e.UrlToParse)
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
				"Bill":      b.Bill,
				"Course":    q.In,
				"MinAmount": strings.Split(q.MinAmount, " ")[0],
				"MaxAmount": strings.Split(q.MaxAmount, " ")[0],
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

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "На данный момент по этой валюте нет поддерживаемых направлений обмена.")

	if len(coins) > 0 {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Какую валюту хочешь получить?")
	}

	msg.ReplyMarkup = m.kbd.Exchange().ReceiveAsResultOfExchangeList(coins, p["From"].(string))

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
