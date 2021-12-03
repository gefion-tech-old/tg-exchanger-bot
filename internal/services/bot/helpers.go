package bot

import (
	"encoding/json"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
)

// Метод слушатель
func (bot *Bot) HandleMessage(m *nsq.Message) error {
	msgEvent := models.MessageEvent{}

	if err := json.Unmarshal(m.Body, &msgEvent); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(msgEvent.To.ChatID, msgEvent.Message.Text)
	bot.botAPI.Send(msg)
	return nil
}
