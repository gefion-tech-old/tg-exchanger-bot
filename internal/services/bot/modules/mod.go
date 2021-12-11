package modules

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules/bills"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules/exchanges"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotModules struct {
	bAPI *tgbotapi.BotAPI
	sAPI api.ApiI
	kbd  keyboards.KeyboardsI

	exchangesMod exchanges.ModExchangesI
	billMod      bills.ModBillsI
}

type BotModulesI interface {
	Exchange() exchanges.ModExchangesI
	Bill() bills.ModBillsI
}

func InitBotModules(bAPI *tgbotapi.BotAPI, kbd keyboards.KeyboardsI, sAPI api.ApiI) BotModulesI {
	return &BotModules{
		bAPI: bAPI,
		sAPI: sAPI,
		kbd:  kbd,
	}
}

func (m *BotModules) Bill() bills.ModBillsI {
	if m.billMod != nil {
		return m.billMod
	}

	m.billMod = bills.InitModBills(m.bAPI, m.sAPI, m.kbd)
	return m.billMod
}

func (m *BotModules) Exchange() exchanges.ModExchangesI {
	if m.exchangesMod != nil {
		return m.exchangesMod
	}

	m.exchangesMod = exchanges.InitModExchanges(m.bAPI, m.sAPI, m.kbd)
	return m.exchangesMod
}
