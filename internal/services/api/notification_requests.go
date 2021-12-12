package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type NotificationRequests struct {
	url string
}

type NotificationRequestsI interface {
	Create(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitNotificationRequests(u string) NotificationRequestsI {
	return &NotificationRequests{
		url: u,
	}
}

func (r *NotificationRequests) Create(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
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
