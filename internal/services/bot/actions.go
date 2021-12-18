package bot

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasthttp"
)

func (bot *Bot) ActionsHandler(ctx context.Context, update tgbotapi.Update, payload map[string]interface{}) error {
	defer tools.Recovery(bot.logger)

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

	case static.BOT__A__EX__NEW_EXCHAGE:
		switch action.Step {
		case 1:
			if p["CbQ"] == static.BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE {
				rMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID-2)
				bot.botAPI.Send(rMsg)

				return bot.m.Exchange().ReceiveAsResultOfExchange(ctx, update, p)
			}

			return bot.m.Exchange().CreateLinkForPayment(ctx, update, &action)
		default:
			return nil
		}

	default:
		return nil
	}
}

func (bot *Bot) CancelAnyAction(ctx context.Context, update tgbotapi.Update, payload map[string]interface{}) error {
	defer tools.Recovery(bot.logger)

	if err := bot.redis.UserActions().Delete(update.Message.Chat.ID); err != nil {
		return err
	}

	if payload["ActionType"] != nil {
		if int(payload["ActionType"].(float64)) == static.BOT__A__EX__NEW_EXCHAGE {
			action := models.UserAction{}
			if err := mapstructure.Decode(payload, &action); err != nil {
				return err
			}

			// Вызываю через повторитель метод отправки уведомления на сервер
			r := api.Retry(bot.sAPI.Notification().Create, 3, time.Second)
			resp, err := r(ctx, map[string]interface{}{
				"type": action.ActionType,
				"meta_data": map[string]interface{}{
					"exchange_from": action.MetaData["From"],
					"exchange_to":   action.MetaData["To"],
				},
				"user": map[string]interface{}{
					"chat_id":  action.User.ChatID,
					"username": action.User.Username,
				},
			})
			if err != nil {
				return errors.ErrBotServerNoAnswer
			}
			defer fasthttp.ReleaseResponse(resp)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Действие успешно отменено.")
	msg.ReplyMarkup = bot.kbd.Base().BaseStartReplyMarkup()
	bot.botAPI.Send(msg)
	return nil
}
