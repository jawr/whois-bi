package sender

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/jhillyerd/enmime"
)

type Sender struct {
	auth smtp.Auth
}

func NewSender() *Sender {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
	)

	sender := Sender{
		auth: auth,
	}

	return &sender
}

func (s Sender) Send(to, subject, body string) error {
	message := enmime.Builder().
		From(os.Getenv("SMTP_FROM_NAME"), os.Getenv("SMTP_EMAIL")).
		To(to, to).
		Subject(subject).
		HTML([]byte(body))

	// can convert this to use tls
	addr := fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))
	if err := message.Send(addr, s.auth); err != nil {
		return err
	}

	return nil
}
