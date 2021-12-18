package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type NotificationRequests struct {
	url    string
	logger *logrus.Logger
}

type NotificationRequestsI interface {
	Create(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitNotificationRequests(u string, l *logrus.Logger) NotificationRequestsI {
	return &NotificationRequests{
		url:    u,
		logger: l,
	}
}

func (r *NotificationRequests) Create(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	defer tools.Recovery(r.logger)

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetBody(b)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("%s/api/v1/admin/notification", r.url))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
