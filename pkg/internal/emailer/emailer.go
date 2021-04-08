package emailer

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	lru "github.com/hashicorp/golang-lru"
	"github.com/jhillyerd/enmime"
	"github.com/segmentio/fasthash/fnv1a"
)

// Emailer sends emails out using the provided sender and using
// the provided email / from name. It uses an LRU cache to prevent
// duplicates being sent out, this can probably be removed once acceptable
// code coverage appears over the codebase.
type Emailer struct {
	// runtime required options
	fromName  string
	fromEmail string

	sender enmime.Sender

	// used to prevent duplicates being sent
	cache *lru.Cache
}

// Creates an SMTP sender using environment variables
func NewSMTPSenderFromEnv() enmime.Sender {
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))
	addr := fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))
	return enmime.NewSMTP(addr, auth)
}

// Create a new Emailer that sends via the provided sender
func NewEmailer(fromName, fromEmail string, sender enmime.Sender) (*Emailer, error) {
	cache, err := lru.New(256)
	if err != nil {
		return nil, err
	}

	emailer := Emailer{
		fromName:  fromName,
		fromEmail: fromEmail,
		cache:     cache,
		sender:    sender,
	}

	return &emailer, nil
}

// Send an email. Hash of the entire email is taken and cached to prevent duplicates
func (s *Emailer) Send(to, subject, body string) error {
	hash := fnv1a.HashString64(to + subject + body)
	if s.cache.Contains(hash) {
		return nil
	}

	s.cache.Add(hash, struct{}{})

	// crude validation
	if len(to) == 0 || !strings.Contains(to, "@") {
		return errors.New("invalid recipient")
	}

	msg := enmime.Builder().
		From(s.fromName, s.fromEmail).
		To(to, to).
		Subject(subject).
		HTML([]byte(body))

	if err := msg.Send(s.sender); err != nil {
		return err
	}

	return nil
}
