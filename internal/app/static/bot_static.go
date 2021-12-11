package static

// Набор поддерживаемых ботом команд
const (
	BOT__CMD__START = "/start"
)

// Набор ресурсов для кнопок
const (
	BOT__BTN__BASE__NEW_EXCHANGE = "💵 Новый обмен"
	BOT__BTN__BASE__MY_BILLS     = "💳 Мои счета"
	BOT__BTN__BASE__MY_EXCHANGES = "📜 Мои обмены"
	BOT__BTN__BASE__SUPPORT      = "🔔 Поддержка"
	BOT__BTN__BASE__ABOUT_BOT    = "ℹ️ О боте"
	BOT__BTN__BASE__OPERATORS    = "🤖 Операторы"
)

// Наборы для CallbackQuery

// Нобор CallbackQuery для модуля Exchanges
const (
	BOT__CQ__EX__COINS_TO_EXCHAGE       = "c_to_ex"      // Показать список валют которые можно обменять
	BOT__CQ__EX__SELECT_COIN_TO_EXCHAGE = "slct_c_to_ex" // Показать список монет которые можно получить
	BOT__CQ__EX__REQ_AMOUNT             = "req_em_ex"    // Запросить сумму обмена
)

// Нобор CallbackQuery для модуля Bills
const (
	BOT__CQ_BL__ADD_BILL = "add_bl_bl"
)
