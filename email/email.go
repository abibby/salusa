package email

import (
	"time"

	"github.com/go-mail/mail"
)

type Mailer interface {
	Mail(*Message) error
}

type Message struct {
	To       []string
	Subject  string
	HTMLBody string
}

type SMTPMailer struct {
	d    *mail.Dialer
	from string
}

var _ Mailer = (*SMTPMailer)(nil)

func NewSMTPMailer(host string, port int, username, password, from string) *SMTPMailer {
	d := mail.NewDialer(host, port, username, password)
	d.Timeout = time.Minute
	return &SMTPMailer{
		d:    d,
		from: from,
	}
}

func (s *SMTPMailer) Mail(m *Message) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", m.To...)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.HTMLBody)

	return s.d.DialAndSend(msg)
}
