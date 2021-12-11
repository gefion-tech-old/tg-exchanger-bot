package keyboards

import (
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var _ BillKeyboardsI = (*BillKeyboards)(nil)

type BillKeyboards struct{}

type BillKeyboardsI interface {
	MyBillsList(arr []models.Bill) tgbotapi.InlineKeyboardMarkup
}

func (kb *BillKeyboards) MyBillsList(arr []models.Bill) tgbotapi.InlineKeyboardMarkup {
	var k tgbotapi.InlineKeyboardMarkup
	for i := 0; i < len(arr); i++ {
		k.InlineKeyboard = append(k.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				arr[i].Bill,
				"",
			),
		))
	}

	k.InlineKeyboard = append(k.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❇️ Добавить ❇️", fmt.Sprintf(`{"CbQ": "%s"}`, static.BOT__CQ_BL__ADD_BILL)),
		))

	return k
}
