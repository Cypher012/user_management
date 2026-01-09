package email

type Sender interface {
	Send(to, subject, html string) error
}
