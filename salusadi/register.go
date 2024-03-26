package salusadi

import (
	"context"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/request"
)

func Register[T auth.User](ctx context.Context) error {
	clog.Register(ctx, nil)
	request.Register(ctx)
	auth.Register[T](ctx)
	return nil
}
