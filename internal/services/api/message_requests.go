package api

import (
	"context"

	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type MessageRequests struct {
	url    string
	logger *logrus.Logger
}

type MessageRequestsI interface {
	Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitMessageRequests(u string, l *logrus.Logger) MessageRequestsI {
	return &MessageRequests{
		url:    u,
		logger: l,
	}
}

func (r *MessageRequests) Get(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	defer tools.Recovery(r.logger)

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
