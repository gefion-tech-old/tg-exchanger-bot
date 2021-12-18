package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UserRequests struct {
	url    string
	logger *logrus.Logger
}

/*
	Сботка методов для работы с пользователскими данными и аккаунтом
*/
type UserRequestsI interface {
	Registration(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitUserRequests(u string, l *logrus.Logger) UserRequestsI {
	return &UserRequests{
		url:    u,
		logger: l,
	}
}

/*
	Регистрация пользователя
*/
func (r *UserRequests) Registration(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	defer tools.Recovery(r.logger)

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetBody(b)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("%s/api/v1/bot/registration", r.url))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
