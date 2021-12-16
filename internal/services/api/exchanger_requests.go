package api

import (
	"context"

	"github.com/valyala/fasthttp"
)

type ExchangerRequests struct {
	url string
}

type ExchangerRequestsI interface {
	Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
	GetQuotesXML(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitExchangerRequests(u string) ExchangerRequestsI {
	return &ExchangerRequests{
		url: u,
	}
}

func (r *ExchangerRequests) Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
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
