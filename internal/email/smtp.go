package email

import (
	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	dialer *gomail.Dialer
	from   string
}

func NewSTMPMailer(host string, port int, username, pass, from string) *SMTPMailer {
	d := gomail.NewDialer(host, port, username, pass)
	return &SMTPMailer{
		dialer: d,
		from:   from,
	}
}

func (m *SMTPMailer) Send(to, subject, html string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", html)

	return m.dialer.DialAndSend(msg)
}
