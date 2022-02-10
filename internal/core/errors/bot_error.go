package errors

import "errors"

var (
	ErrBotServerNoAnswer = errors.New("не удалось выплонить запрос, сервер не отвечает")
)
