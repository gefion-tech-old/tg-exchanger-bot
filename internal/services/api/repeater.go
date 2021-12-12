package api

import (
	"context"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type UserReqStructCtx int8
type MessageConnectorCtx string

const UserReqStructCtxKey UserReqStructCtx = iota
const MessageConnectorCtxKey MessageConnectorCtx = "connector"

// Сигнатура функции взаимодействующая со службой
type Effector func(context.Context, map[string]interface{}) (*fasthttp.Response, error)

//	Учитывает возможный временный характер ошибки  и
//	осуществляет повторные попытки выполнить неучаную операцию.
func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context, body map[string]interface{}) (*fasthttp.Response, error) {
		for r := 0; ; r++ {
			resp, err := effector(ctx, body)
			if err == nil || r >= retries {
				return resp, err
			}

			fmt.Printf("Attempt %d failed; retrying in %v\n", r+1, delay)

			delay += time.Second

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}
