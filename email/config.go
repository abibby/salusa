package email

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
