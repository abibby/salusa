package providers

import (
	"context"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/static/template/config"
)

// var ModelRegistrar = modeldi.NewModelRegistrar()
var registrar = []func(*di.DependencyProvider){}

func Add(register func(*di.DependencyProvider)) {
	registrar = append(registrar, register)
}

// Register registers any custom di providers
func Register(dp *di.DependencyProvider) {
	for _, register := range registrar {
		register(dp)
	}

	di.RegisterWith(dp, func(ctx context.Context, tag string, cfg *config.Config) (email.Mailer, error) {
		return cfg.Mail.Mailer(), nil
	})
}
