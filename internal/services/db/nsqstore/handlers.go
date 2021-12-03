package nsqstore

import (
	"encoding/json"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
)

type EventsHandler struct {
	BotAPI *tgbotapi.BotAPI
}

func (e *EventsHandler) HandleMessage(m *nsq.Message) error {
	msgEvent := models.MessageEvent{}

	if err := json.Unmarshal(m.Body, &msgEvent); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(msgEvent.To.ChatID, msgEvent.Message.Text)
	e.BotAPI.Send(msg)
	return nil
}
