package salusadi

import (
	"context"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/filesystem"
	"github.com/abibby/salusa/openapidoc/openapidocdi"
	"github.com/abibby/salusa/request"
)

func Register[TUser auth.User](migrations *migrate.Migrations) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		registerers := []func(context.Context) error{
			clog.Register,
			request.Register,
			auth.Register[TUser],
			databasedi.RegisterFromConfig(migrations),
			email.Register,
			event.Register,
			filesystem.Register,
			openapidocdi.Register,
		}

		for _, register := range registerers {
			err := register(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
