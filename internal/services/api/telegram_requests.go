package api

import (
	"context"
	"fmt"

	"github.com/gefion-tech/tg-exchanger-bot/internal/tools"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type TelegramRequests struct {
	url    string
	token  string
	logger *logrus.Logger
}

type TelegramRequestsI interface {
	GetFileDate(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
	DownloadFile(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error)
}

func InitTelegramRequests(u, t string, l *logrus.Logger) TelegramRequestsI {
	return &TelegramRequests{
		url:    u,
		token:  t,
		logger: l,
	}

}

func (r *TelegramRequests) GetFileDate(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	defer tools.Recovery(r.logger)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("%sbot%s/getFile?file_id=%s", r.url, r.token, body["file_id"]))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}

func (r *TelegramRequests) DownloadFile(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("%sfile/bot%s/%s", r.url, r.token, body["file_path"]))
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return nil, err
	}

	defer fasthttp.ReleaseRequest(req)
	return res, nil
}
