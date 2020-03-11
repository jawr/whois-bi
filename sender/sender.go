package sender

import (
	"net/smtp"
	"os"

	"github.com/jhillyerd/enmime"
)

const (
	SmtpName     string = "Whois.bi"
	SmtpUsername string = "hi@whois.bi"
	SmtpAddr     string = "arrow.mxrouting.net"
	SmtpPort     string = "25"
)

type Sender struct {
	auth smtp.Auth
}

func NewSender() *Sender {
	auth := smtp.PlainAuth(
		"",
		SmtpUsername,
		os.Getenv("MONERE_SMTP_PASSWORD"),
		SmtpAddr,
	)

	sender := Sender{
		auth: auth,
	}

	return &sender
}

func (s Sender) Send(to, subject, body string) error {
	message := enmime.Builder().
		From(SmtpName, SmtpUsername).
		To(to, to).
		Subject(subject).
		HTML([]byte(body))

		// can convert this to use tls
	if err := message.Send(SmtpAddr+":"+SmtpPort, s.auth); err != nil {
		return err
	}

	return nil
}
