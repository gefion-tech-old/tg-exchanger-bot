package keyboards

type Keyboards struct {
	baseKeyboards     *BaseKeyboards
	exchangeKeyboards *ExchangeKeyboards
}

type KeyboardsI interface {
	Base() BaseKeyboardsI
	Exchange() ExchangeKeyboardsI
}

func InitKeyboards() KeyboardsI {
	return &Keyboards{}
}

func (c *Keyboards) Exchange() ExchangeKeyboardsI {
	if c.exchangeKeyboards != nil {
		return c.exchangeKeyboards
	}

	c.exchangeKeyboards = &ExchangeKeyboards{}
	return c.exchangeKeyboards
}

func (c *Keyboards) Base() BaseKeyboardsI {
	if c.baseKeyboards != nil {
		return c.baseKeyboards
	}

	c.baseKeyboards = &BaseKeyboards{}
	return c.baseKeyboards
}
