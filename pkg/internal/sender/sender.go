package sender

import (
	"fmt"
	"net/smtp"
	"os"

	lru "github.com/hashicorp/golang-lru"
	"github.com/jhillyerd/enmime"
	"github.com/segmentio/fasthash/fnv1a"
)

type Sender struct {
	auth  smtp.Auth
	cache *lru.Cache
}

func NewSender() (*Sender, error) {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
	)

	cache, err := lru.New(256)
	if err != nil {
		return nil, err
	}

	sender := Sender{
		auth:  auth,
		cache: cache,
	}

	return &sender, nil
}

func (s *Sender) Send(to, subject, body string) error {
	hash := fnv1a.HashString64(to + subject + body)
	if s.cache.Contains(hash) {
		return nil
	}

	s.cache.Add(hash, struct{}{})

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
