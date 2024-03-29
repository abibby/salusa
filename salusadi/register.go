package salusadi

import (
	"context"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/request"
)

func Register[T auth.User](ctx context.Context) error {
	clog.Register(ctx, nil)
	err := request.Register(ctx)
	if err != nil {
		return err
	}
	err = auth.Register[T](ctx)
	if err != nil {
		return err
	}
	return nil
}
