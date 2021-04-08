package emailer

import (
	"bytes"
	"errors"
	"sync"

	"github.com/jhillyerd/enmime"
)

type email struct {
	from string
	to   []string
	body []byte
}

// satisfies the Sender interface but keeps all messages in memory
type MemorySender struct {
	emails []email
	sync.Mutex
	err error
}

// create a new MemorySender
func NewMemorySender() *MemorySender {
	return &MemorySender{
		emails: make([]email, 0),
	}
}

// append the message in to memory
func (s *MemorySender) Send(from string, to []string, msg []byte) error {
	s.Lock()
	defer s.Unlock()
	if s.err != nil {
		err := s.err
		s.err = nil
		return err
	}
	s.emails = append(s.emails, email{from, to, msg})
	return nil
}

// retrieve a message sent
func (s *MemorySender) EmailAt(idx int) (*enmime.Envelope, error) {
	s.Lock()
	defer s.Unlock()

	if len(s.emails) == 0 || len(s.emails) < idx || idx < 0 {
		return nil, errors.New("bad index")
	}

	return enmime.ReadEnvelope(bytes.NewReader(s.emails[idx].body))
}
