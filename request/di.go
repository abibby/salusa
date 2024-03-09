package request

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abibby/salusa/di"
)

type contextKey uint8

const (
	requestKey contextKey = iota
	responseKey
)

func RegisterDI(dp *di.DependencyProvider) error {
	di.Register(dp, func(ctx context.Context, tag string) (*http.Request, error) {
		req, ok := ctx.Value(requestKey).(*http.Request)
		if !ok {
			return nil, fmt.Errorf("request not in context")
		}
		return req, nil
	})
	di.Register(dp, func(ctx context.Context, tag string) (http.ResponseWriter, error) {
		resp, ok := ctx.Value(responseKey).(http.ResponseWriter)
		if !ok {
			return nil, fmt.Errorf("response not in context")
		}
		return resp, nil
	})
	return nil
}
