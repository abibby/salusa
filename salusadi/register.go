package salusadi

import (
	"context"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
)

type Config interface {
	dialects.DBConfiger
	kernel.KernelConfig
}

func Register[TConfig Config, TUser auth.User](migrations *migrate.Migrations) func(ctx context.Context) error {
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
		err = databasedi.RegisterFromConfig[TConfig](migrations)(ctx)
		if err != nil {
			return err
		}
		err = email.RegisterMailer[TConfig](ctx)
		if err != nil {
			return err
		}
		err = event.Register[TConfig](ctx)
		if err != nil {
			return err
		}
		return nil
	}
}
