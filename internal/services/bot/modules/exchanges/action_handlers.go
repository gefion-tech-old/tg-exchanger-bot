package exchanges

import (
	"context"
	"regexp"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (m *ModExchanges) CreateLinkForPayment(ctx context.Context, update tgbotapi.Update, action *models.UserAction) error {
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

		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хм, это непохоже на сумму...")
	m.bAPI.Send(msg)
	return nil
}
