package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func GetMessage(ctx context.Context, update tgbotapi.Update, sAPI api.ApiI, connector string, logger *logrus.Logger, params ...interface{}) (*models.Message, error) {
	defer tools.Recovery(logger)

	ctx = context.WithValue(ctx, api.MessageConnectorCtxKey, connector)
	r := api.Retry(sAPI.Message().Get, 3, time.Second)

	resp, err := r(ctx, map[string]interface{}{"connector": connector})
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	msg := models.Message{}
	if err := json.Unmarshal(resp.Body(), &msg); err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		msg.MessageText = fmt.Sprintf(msg.MessageText, params)
		return &msg, nil
	case http.StatusNotFound:
		msg.MessageText = "Сообщение не найдено..."
		return &msg, nil
	default:
		msg.MessageText = "Произошла ошибка при получении этого сообщения, попробуйте повторить попытку позже."
		return &msg, nil
	}
}
