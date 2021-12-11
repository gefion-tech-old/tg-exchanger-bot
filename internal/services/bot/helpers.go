package bot

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
)

// Метод слушатель NSQ событий
func (bot *Bot) HandleMessage(m *nsq.Message) error {
	msgEvent := models.MessageEvent{}
	if err := json.Unmarshal(m.Body, &msgEvent); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(msgEvent.To.ChatID, msgEvent.Message.Text)
	bot.botAPI.Send(msg)
	return nil
}

// Метод хелпер для обработки ошибок из нижестоящих хендлеров
func (bot *Bot) error(update tgbotapi.Update, errs ...error) {
	for i := 0; i < len(errs); i++ {
		if errs[i] != nil {
			msg := tgbotapi.NewMessage(int64(update.Message.Chat.ID), "❌ Похоже произошла какая-то ошибка ❌")
			if errs[i] == sql.ErrNoRows {
				msg = tgbotapi.NewMessage(int64(update.Message.Chat.ID), "❌ Запрашиваемый ресурс не найден ❌")
			}
			bot.botAPI.Send(msg)

			// Отправка отладочной информации об ошибке если пользователь разработчик
			for i := 0; i < len(bot.cnf.Developers); i++ {
				if update.Message.Chat.UserName == bot.cnf.Developers[i] {
					errMsg := tgbotapi.NewMessage(int64(update.Message.Chat.ID), fmt.Sprintf("`%s`", errs[i].Error()))
					errMsg.ParseMode = tgbotapi.ModeMarkdown
					bot.botAPI.Send(errMsg)
					break
				}
			}
		}
	}
}

// Метод проверки пользовательского username
// Если username отсутствует, бот не будет дальше пускать
// username необходим для корректной работы при выполнение последующих операций
func (bot *Bot) check(update tgbotapi.Update) bool {
	if update.Message.Chat.UserName == "" {
		msg := tgbotapi.NewMessage(int64(update.Message.Chat.ID), "Привет! Вам необходимо установить себе `username` для использования этого бота. После этого выполните команду /start чтобы перезапустить бота.")
		msg.ParseMode = tgbotapi.ModeMarkdown
		bot.botAPI.Send(msg)
		return false
	}
	return true
}

func (bot *Bot) rewriter(update tgbotapi.Update) tgbotapi.Update {
	// Если получили Message а CallbackQuery пустой
	if update.Message != nil {
		update.CallbackQuery = &tgbotapi.CallbackQuery{
			Message: &tgbotapi.Message{
				MessageID: update.Message.MessageID,
				Chat: &tgbotapi.Chat{
					ID:        update.Message.Chat.ID,
					UserName:  update.Message.Chat.UserName,
					FirstName: update.Message.Chat.FirstName,
					LastName:  update.Message.Chat.LastName,
				},
			},
		}
		return update
	}

	// Если получили CallbackQuery а Message пустой
	update.Message = &tgbotapi.Message{
		MessageID: update.CallbackQuery.Message.MessageID,
		Chat: &tgbotapi.Chat{
			ID:        update.CallbackQuery.Message.Chat.ID,
			UserName:  update.CallbackQuery.Message.Chat.UserName,
			FirstName: update.CallbackQuery.Message.Chat.FirstName,
			LastName:  update.CallbackQuery.Message.Chat.LastName,
		},
	}
	return update
}
