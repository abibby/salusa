package salusadi

import (
	"context"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
)

func Register[TUser auth.User](migrations *migrate.Migrations) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		clog.Register(ctx, nil)
		err := request.Register(ctx)
		if err != nil {
			return err
		}
		err = auth.Register[TUser](ctx)
		if err != nil {
			return err
		}
		err = databasedi.RegisterFromConfig[kernel.KernelConfig](migrations)(ctx)
		if err != nil {
			return err
		}
		err = email.RegisterMailer[kernel.KernelConfig](ctx)
		if err != nil {
			return err
		}
		err = event.Register[kernel.KernelConfig](ctx)
		if err != nil {
			return err
		}
		return nil
	}
}
