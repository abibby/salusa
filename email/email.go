package email

import (
	"github.com/go-mail/mail"
)

type Mailer interface {
	Mail(*Message) error
}

type Message struct {
	From     string
	To       []string
	Subject  string
	HTMLBody string
}

type SMTPMailer struct {
	d *mail.Dialer
}

var _ Mailer = (*SMTPMailer)(nil)

func NewSMTPMailer(host string, port int, username, password string) *SMTPMailer {
	return &SMTPMailer{
		d: mail.NewDialer(host, port, username, password),
	}
}

func (s *SMTPMailer) Mail(m *Message) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To...)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.HTMLBody)

	return s.d.DialAndSend(msg)
}
