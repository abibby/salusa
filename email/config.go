package email

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
)

type MailConfiger interface {
	MailConfig() Config
}
type Config interface {
	Mailer() Mailer
}

type SMTPConfig struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func (c *SMTPConfig) Mailer() Mailer {
	return NewSMTPMailer(c.Host, c.Port, c.Username, c.Password, c.From)
}

func RegisterMailer[T kernel.KernelConfig](ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(cfg T) (Mailer, error) {
		var cfgAny any = cfg
		cfger, ok := cfgAny.(MailConfiger)
		if !ok {
			return nil, fmt.Errorf("config not instance of email.MailConfiger")
		}
		return cfger.MailConfig().Mailer(), nil
	})
	return nil
}
