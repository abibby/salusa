package providers

import (
	"context"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/static/template/config"
)

// var ModelRegistrar = modeldi.NewModelRegistrar()
var registrar = []func(context.Context){}

func Add(register func(context.Context)) {
	registrar = append(registrar, register)
}

// Register registers any custom di providers
func Register(ctx context.Context) error {
	for _, register := range registrar {
		register(ctx)
	}

	di.RegisterWith(ctx, func(ctx context.Context, tag string, cfg *config.Config) (email.Mailer, error) {
		return cfg.Mail.Mailer(), nil
	})
	return nil
}
