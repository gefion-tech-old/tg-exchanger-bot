package api

import (
	"context"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/models"
	"github.com/valyala/fasthttp"
)

type BillRequests struct {
	url string
}

type BillRequestsI interface {
	GetAll(ctx context.Context) (*fasthttp.Response, error)
}

func InitBillRequests(u string) BillRequestsI {
	return &BillRequests{
		url: u,
	}
}

func (r *BillRequests) GetAll(ctx context.Context) (*fasthttp.Response, error) {
	ctxUserReq := ctx.Value(UserReqStructCtxKey).(*models.UserReq)
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("%s/api/v1/bot/user/%d/bills", r.url, ctxUserReq.ChatID))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
