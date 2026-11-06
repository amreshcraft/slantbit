package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host, port, user, pass, from string
}

func New(host, port, user, pass, from string) *Mailer {
	return &Mailer{host, port, user, pass, from}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body)

	auth := smtp.PlainAuth("", m.user, m.pass, m.host)
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{to}, []byte(msg))
}