package exchanges

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var _ ModExchangesI = (*ModExchanges)(nil)

type ModExchanges struct {
	bAPI   *tgbotapi.BotAPI
	sAPI   api.ApiI
	redis  redisstore.RedisStoreI
	kbd    keyboards.KeyboardsI
	logger *logrus.Logger
}

type ModExchangesI interface {
	/* Action обработчики */

	CreateLinkForPayment(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error

	HandleReceivedAddress(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error

	/* CallbackQuery обработчики */

	// @CallbackQuery BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE
	ReceiveAsResultOfExchange(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error

	СhooseBill(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error

	// @CallbackQuery BOT__CQ__EX__REQ_AMOUNT
	ReqAmount(ctx context.Context, update tgbotapi.Update, p map[string]interface{}) error

	/*  Универсальные обработчики */

	// @Button BOT__BTN__BASE__NEW_EXCHANGE
	// @CallbackQuery BOT__CQ__EX__COINS_TO_EXCHAGE
	NewExchange(ctx context.Context, update tgbotapi.Update) error
}

func InitModExchanges(bAPI *tgbotapi.BotAPI, servAPI api.ApiI, redis redisstore.RedisStoreI, k keyboards.KeyboardsI, l *logrus.Logger) ModExchangesI {
	return &ModExchanges{
		bAPI:   bAPI,
		sAPI:   servAPI,
		redis:  redis,
		kbd:    k,
		logger: l,
	}
}

func (m *ModExchanges) GetExchangeInfo(ctx context.Context, update tgbotapi.Update, from, to string) (*models.OneObmenItem, error) {
	// Получение обменника
	r := api.Retry(m.sAPI.Exchanger().Get, 3, time.Second)
	resp, err := r(ctx, map[string]interface{}{
		"name": "1obmen",
	})
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер не отвечает")
		m.bAPI.Send(msg)
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != http.StatusOK {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Не удалось получить данные актуального курса ❌")
		m.bAPI.Send(msg)
		return nil, errors.New("не удалось получить данные актуального курса")
	}

	body := models.Exchanger{}
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}

	// Получение актуальных котировок
	data, err := m.quotes(ctx, update, body.UrlToParse)
	if err != nil {
		return nil, err
	}

	// Поиск нужной котировки
	cQ := make(chan *models.OneObmenItem)
	for i := 0; i < len(data.Rates); i++ {
		go func(i int) {
			if data.Rates[i].From == from && data.Rates[i].To == to {
				defer close(cQ)
				cQ <- &data.Rates[i]
			}
		}(i)
	}

	q := <-cQ

	return q, nil
}
