package api

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
)

type MessageRequests struct {
	url string
}

type MessageRequestsI interface {
	Get(ctx context.Context) (*fasthttp.Response, error)
}

func InitMessageRequests(u string) MessageRequestsI {
	return &MessageRequests{
		url: u,
	}
}

func (r *MessageRequests) Get(ctx context.Context) (*fasthttp.Response, error) {
	connector := ctx.Value(MessageConnectorCtxKey).(string)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("%s/api/v1/admin/message?connector=%s", r.url, connector))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
