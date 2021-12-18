package api

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ExchangerRequests struct {
	url    string
	logger *logrus.Logger
}

type ExchangerRequestsI interface {
	Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
	GetQuotesXML(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitExchangerRequests(u string, l *logrus.Logger) ExchangerRequestsI {
	return &ExchangerRequests{
		url:    u,
		logger: l,
	}
}

func (r *ExchangerRequests) Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	defer tools.Recovery(r.logger)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(r.url + "/api/v1/admin/exchanger/" + body["name"].(string))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}

func (r *ExchangerRequests) GetQuotesXML(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.SetRequestURI(body["url"].(string))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
