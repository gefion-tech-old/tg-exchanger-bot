package bot

import (
	"context"
	"encoding/json"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mitchellh/mapstructure"
)

func (bot *Bot) ActionsHandler(ctx context.Context, update tgbotapi.Update, payload map[string]interface{}) error {
	action := models.UserAction{}
	if err := mapstructure.Decode(payload, &action); err != nil {
		return err
	}

	p := map[string]interface{}{}
	if update.CallbackQuery != nil {
		bot.error(bot.rewriter(update), json.Unmarshal([]byte(update.CallbackQuery.Data), &p))
	}

	switch action.ActionType {
	case static.BOT__A__BL__ADD_NEW_BILL:
		switch action.Step {
		case 1:
			return bot.m.Bill().AddNewBillStepTwo(ctx, update, &action)

		case 2:
			if p["CbQ"] == static.BOT__CQ_BL__ADD_BILL_VALID_S_2 {
				return bot.m.Bill().AddNewBillStepThreeCorrect(ctx, bot.rewriter(update), &action)
			}

			if p["CbQ"] == static.BOT__CQ_BL__ADD_BILL_N_VALID_S_2 {
				return bot.m.Bill().AddNewBillStepThreeInCorrect(ctx, bot.rewriter(update), &action)
			}

			return nil

		case 3:
			return bot.m.Bill().AddNewBillStepFour(ctx, update, &action)

		default:
			return nil
		}

	default:
		return nil
	}
}

func (bot *Bot) CancelAnyAction(update tgbotapi.Update) error {
	if err := bot.redis.UserActions().Delete(update.Message.Chat.ID); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Действие успешно отменено.")
	msg.ReplyMarkup = bot.kbd.Base().BaseStartReplyMarkup()
	bot.botAPI.Send(msg)
	return nil
}
