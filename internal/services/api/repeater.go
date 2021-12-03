package api

import (
	"context"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type UserReqStructCtx int8

const UserReqStructCtxKey UserReqStructCtx = iota

// Сигнатура функции взаимодействующая со службой
type Effector func(context.Context) (*fasthttp.Response, error)

/*
	Учитывает возможный временный характер ошибки  и
	осуществляет повторные попытки выполнить неучаную операцию.
*/
func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (*fasthttp.Response, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
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
