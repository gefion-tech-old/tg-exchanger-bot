package static

// Набор поддерживаемых ботом команд
const (
	BOT__CMD__START = "/start"
	BOT__CMD__SKIP  = "/skip"
)

// Набор ресурсов для кнопок
const (
	BOT__BTN__BASE__NEW_EXCHANGE = "💵 Новый обмен"
	BOT__BTN__BASE__MY_BILLS     = "💳 Мои счета"
	BOT__BTN__BASE__MY_EXCHANGES = "📜 Мои обмены"
	BOT__BTN__BASE__SUPPORT      = "🔔 Поддержка"
	BOT__BTN__BASE__ABOUT_BOT    = "ℹ️ О боте"
	BOT__BTN__BASE__OPERATORS    = "🤖 Операторы"

	BOT__BTN__OP__CANCEL = "❌ Отменить"
	BOT__BTN__OP__SAVE   = "✅ Завершить"
)

/* Наборы для CallbackQuery */

// Нобор CallbackQuery для модуля Exchanges
const (
	BOT__CQ__EX__COINS_TO_EXCHAGE       = "c_to_ex"      // Показать список валют которые можно обменять
	BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE = "slct_c_to_ex" // Показать список монет которые можно получить
	BOT__CQ__EX__REQ_AMOUNT             = "req_em_ex"    // Запросить сумму обмена
)

// Нобор CallbackQuery для модуля Bills
const (
	BOT__CQ_BL__ADD_BILL_S_1         = "add_bl_s_1"
	BOT__CQ_BL__ADD_BILL_VALID_S_2   = "add_bl_val_s_2"
	BOT__CQ_BL__ADD_BILL_N_VALID_S_2 = "add_bl_nval_s_2"
)

// Набор возможных типов пользовательских действий
const (
	BOT__A__BL__ADD_NEW_BILL = 854
)
