package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/valyala/fasthttp"
)

type UserRequests struct {
	url string
}

/*
	Сботка методов для работы с пользователскими данными и аккаунтом
*/
type UserRequestsI interface {
	Registration(ctx context.Context) (*fasthttp.Response, error)
}

func InitUserRequests(u string) UserRequestsI {
	return &UserRequests{
		url: u,
	}
}

/*
	Регистрация пользователя
*/
func (r *UserRequests) Registration(ctx context.Context) (*fasthttp.Response, error) {
	ctxUserReq := ctx.Value(UserReqStructCtxKey).(*models.UserReq)

	b, err := json.Marshal(map[string]interface{}{
		"chat_id":  ctxUserReq.ChatID,
		"username": ctxUserReq.Username,
	})
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
