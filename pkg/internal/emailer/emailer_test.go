package emailer

import (
	"errors"
	"testing"
)

const (
	fromName  = "mr. tester"
	fromEmail = "test@whois.bi"
	toEmail   = "user@whois.bi"
	subject   = "alert alert alert"
	body      = "hi there!"
)

func Test_Send(t *testing.T) {

	sender := NewMemorySender()

	emailer, err := NewEmailer(fromName, fromEmail, sender)
	if err != nil {
		t.Fatalf("NewEmailer() expected nil got %q", err)
	}

	if err := emailer.Send(toEmail, subject, body); err != nil {
		t.Fatalf("Send() expected nil got %q", err)
	}

	if len(sender.emails) != 1 {
		t.Fatal("Expected sender.emails to be 1")
	}

	env, err := sender.EmailAt(0)
	if err != nil {
		t.Fatalf("EmailAt() expected nil got %q", err)
	}

	if env.GetHeader("Subject") != subject {
		t.Fatal("Unexpected subject")
	}

}

func Test_SendError(t *testing.T) {

	sender := NewMemorySender()

	merr := errors.New("mock error")
	sender.err = merr

	emailer, err := NewEmailer(fromName, fromEmail, sender)
	if err != nil {
		t.Fatalf("NewEmailer() expected nil got %q", err)
	}

	err = emailer.Send(toEmail, subject, body)
	if err == nil {
		t.Fatalf("Send() expected nil got %q", err)
	}

	if err.Error() != merr.Error() {
		t.Fatalf("Send() expected %q got %q", merr, err)
	}
}

func Test_DuplicateSend(t *testing.T) {
	sender := NewMemorySender()

	emailer, err := NewEmailer(fromName, fromEmail, sender)
	if err != nil {
		t.Fatalf("NewEmailer() expected nil got %q", err)
	}

	for i := 0; i < 10; i++ {
		if err := emailer.Send(toEmail, subject, body); err != nil {
			t.Fatalf("Send() expected nil got %q", err)
		}
	}

	if len(sender.emails) != 1 {
		t.Fatal("Expected sender.emails to be 1")
	}

	if err := emailer.Send("another@whois.bi", subject, body); err != nil {
		t.Fatalf("Send() expected nil got %q", err)
	}

	if len(sender.emails) != 2 {
		t.Fatal("Expected sender.emails to be 2")
	}
}

func Test_BadEmail(t *testing.T) {
	sender := NewMemorySender()

	emailer, err := NewEmailer(fromName, fromEmail, sender)
	if err != nil {
		t.Fatalf("NewEmailer() expected nil got %q", err)
	}

	if err := emailer.Send("", subject, body); err == nil {
		t.Fatal("Send() expected an error got nil")
	}
}

func Test_BadMemoryIndex(t *testing.T) {
	sender := NewMemorySender()

	_, err := sender.EmailAt(0)
	if err == nil {
		t.Fatal("EmailAt() expected an error got nil")
	}
}
