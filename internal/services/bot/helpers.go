package bot

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/core/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nsqio/go-nsq"
)

// Метод слушатель NSQ событий
func (bot *Bot) HandleMessage(m *nsq.Message) error {
	defer tools.Recovery(bot.logger)

	nMsg := models.MessageEvent{}
	if err := json.Unmarshal(m.Body, &nMsg); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(nMsg.To.ChatID, nMsg.Message.Text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.botAPI.Send(msg)
	return nil
}

// Игнорировать все команды и кнопки если есть активное действие
// Исключением является команда /skip которая позволяет отменить
// любое начатое действие
func (bot *Bot) action(update tgbotapi.Update) (map[string]interface{}, bool) {
	defer tools.Recovery(bot.logger)

	data, err := bot.redis.UserActions().Fetch(int64(update.Message.Chat.ID))
	switch err {
	// Значит есть активное действие
	case nil:
		ignoreList := []string{
			static.BOT__CMD__SKIP,
			static.BOT__BTN__OP__CANCEL,
		}

		for i := 0; i < len(ignoreList); i++ {
			if update.Message.Text == ignoreList[i] {
				return data, true
			}
		}
		return data, false

	// Активных действий нет
	case redis.Nil:
		return nil, true

	default:
		bot.error(update, err)
		return nil, true
	}
}

// Метод хелпер для обработки ошибок из нижестоящих хендлеров
func (bot *Bot) error(update tgbotapi.Update, errs ...error) {
	defer tools.Recovery(bot.logger)

	for i := 0; i < len(errs); i++ {
		if errs[i] != nil {
			fmt.Println(errs[i].Error())
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

// Вспомогательный метод
// Для перезаписи объекта update
func (bot *Bot) rewriter(update tgbotapi.Update) tgbotapi.Update {
	defer tools.Recovery(bot.logger)

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
