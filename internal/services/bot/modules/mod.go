package modules

import (
	"github.com/gefion-tech/tg-exchanger-bot/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/api"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/keyboards"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules/base"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules/bills"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/bot/modules/exchanges"
	"github.com/gefion-tech/tg-exchanger-bot/internal/services/db/redisstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotModules struct {
	bAPI  *tgbotapi.BotAPI
	sAPI  api.ApiI
	cnf   *config.BotConfig
	redis redisstore.RedisStoreI
	kbd   keyboards.KeyboardsI

	exchangesMod exchanges.ModExchangesI
	billMod      bills.ModBillsI
	baseMod      base.ModBaseI
}

type BotModulesI interface {
	Base() base.ModBaseI
	Exchange() exchanges.ModExchangesI
	Bill() bills.ModBillsI
}

func InitBotModules(bAPI *tgbotapi.BotAPI, kbd keyboards.KeyboardsI, redis redisstore.RedisStoreI, sAPI api.ApiI, cnf *config.BotConfig) BotModulesI {
	return &BotModules{
		bAPI:  bAPI,
		sAPI:  sAPI,
		redis: redis,
		kbd:   kbd,

		cnf: cnf,
	}
}

func (m *BotModules) Base() base.ModBaseI {
	if m.baseMod != nil {
		return m.baseMod
	}

	m.baseMod = base.InitModBase(m.bAPI, m.sAPI, m.redis, m.kbd)
	return m.baseMod
}

func (m *BotModules) Bill() bills.ModBillsI {
	if m.billMod != nil {
		return m.billMod
	}

	m.billMod = bills.InitModBills(m.bAPI, m.sAPI, m.redis, m.kbd, m.cnf)
	return m.billMod
}

func (m *BotModules) Exchange() exchanges.ModExchangesI {
	if m.exchangesMod != nil {
		return m.exchangesMod
	}

	m.exchangesMod = exchanges.InitModExchanges(m.bAPI, m.sAPI, m.redis, m.kbd)
	return m.exchangesMod
}
