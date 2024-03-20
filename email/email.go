package email

import (
	"net/smtp"
	"strings"
)

type Mailer interface {
	Mail(*Message) error
}

type Message struct {
	// From string
	To   []string
	Body []byte
}

type SMTPMailer struct {
	from    string
	address string
	auth    smtp.Auth
}

var _ Mailer = (*SMTPMailer)(nil)

func NewSMTPMailer(from, address string, auth smtp.Auth) *SMTPMailer {
	return &SMTPMailer{
		from:    from,
		address: address,
		auth:    auth,
	}
}

func NewSMTPMailerPlainAuth(from, address, username, password string) *SMTPMailer {
	host := strings.SplitN(address, ":", 2)[0]
	return &SMTPMailer{
		from:    from,
		address: address,
		auth:    smtp.PlainAuth("", username, password, host),
	}
}

func (s *SMTPMailer) Mail(m *Message) error {
	return smtp.SendMail(s.address, s.auth, s.from, m.To, m.Body)
}
