package providers

import (
	"context"

	"github.com/abibby/salusa/database/model/modeldi"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/static/template/config"
)

// Register registers any custom di providers
func Register(dp *di.DependencyProvider) {
	modeldi.Register[*models.User](dp)

	di.Register(dp, func(ctx context.Context, tag string) (email.Mailer, error) {
		cfg, err := di.Resolve[*config.Config](ctx, dp)
		if err != nil {
			return nil, err
		}
		return cfg.Mail.Mailer(), nil
	})
}
