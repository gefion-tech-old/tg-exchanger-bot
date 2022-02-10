package keyboards

import (
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/core/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ BillKeyboardsI = (*BillKeyboards)(nil)

type BillKeyboards struct{}

type BillKeyboardsI interface {
	// InlineKeyboards
	MyBillsList(arr []models.Bill) tgbotapi.InlineKeyboardMarkup
	CardСorrectnessConfirmation() tgbotapi.InlineKeyboardMarkup
}

func (kb *BillKeyboards) CardСorrectnessConfirmation() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Нет", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ_BL__ADD_BILL_N_VALID_S_2)),
			tgbotapi.NewInlineKeyboardButtonData("Да", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ_BL__ADD_BILL_VALID_S_2)),
		),
	)
}

func (kb *BillKeyboards) MyBillsList(arr []models.Bill) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	for i := 0; i < len(arr); i++ {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Bill,
				fmt.Sprintf(`{"CbQ": "%s", "ID": %d}`, static.BOT__CQ_BL__SELECT_BILL, arr[i].ID),
			),
		))
	}

	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ_BL__ADD_BILL_S_1)),
		))

	return k
}
