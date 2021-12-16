package api

import (
	"context"

	"github.com/valyala/fasthttp"
)

type MessageRequests struct {
	url string
}

type MessageRequestsI interface {
	Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitMessageRequests(u string) MessageRequestsI {
	return &MessageRequests{
		url: u,
	}
}

func (r *MessageRequests) Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(r.url + "/api/v1/admin/message/" + body["connector"].(string))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
