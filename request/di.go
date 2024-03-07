package request

import (
	"context"
	"log"
	"net/http"
	"reflect"

	"github.com/abibby/salusa/di"
)

type contextKey uint8

const (
	requestKey contextKey = iota
	responseKey
)

func InitDI(context.Context) error {
	di.Register(func(ctx context.Context, tag string) *http.Request {

		v := ctx.Value(requestKey)
		req, ok := ctx.Value(requestKey).(*http.Request)
		log.Printf("%#v", reflect.TypeOf(v).String())
		if !ok {
			return nil
		}
		return req
	})
	di.Register(func(ctx context.Context, tag string) http.ResponseWriter {
		resp, ok := ctx.Value(responseKey).(http.ResponseWriter)
		if !ok {
			return nil
		}
		return resp
	})
	return nil
}
