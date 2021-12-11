package keyboards

type Keyboards struct {
	baseKeyboards     BaseKeyboardsI
	exchangeKeyboards ExchangeKeyboardsI
	billKeyboards     BillKeyboardsI
}

type KeyboardsI interface {
	Base() BaseKeyboardsI
	Exchange() ExchangeKeyboardsI
	Bill() BillKeyboardsI
}

func InitKeyboards() KeyboardsI {
	return &Keyboards{}
}

func (k *Keyboards) Bill() BillKeyboardsI {
	if k.billKeyboards != nil {
		return k.billKeyboards
	}

	k.billKeyboards = &BillKeyboards{}
	return k.billKeyboards
}

func (k *Keyboards) Exchange() ExchangeKeyboardsI {
	if k.exchangeKeyboards != nil {
		return k.exchangeKeyboards
	}

	k.exchangeKeyboards = &ExchangeKeyboards{}
	return k.exchangeKeyboards
}

func (k *Keyboards) Base() BaseKeyboardsI {
	if k.baseKeyboards != nil {
		return k.baseKeyboards
	}

	k.baseKeyboards = &BaseKeyboards{}
	return k.baseKeyboards
}
