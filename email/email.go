package email

import "net/smtp"

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

func (s *SMTPMailer) Mail(m *Message) error {
	return smtp.SendMail(s.address, s.auth, s.from, m.To, m.Body)
}
