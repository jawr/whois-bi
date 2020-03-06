package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jawr/monere/sender"
	"github.com/pkg/errors"
)

var (
	email = flag.String("email", "", "email to add")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -email\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	s := sender.NewSender()

	if err := s.Send("jess@lawrence.pm", "my subject", "Hey"); err != nil {
		return errors.Wrap(err, "Send")
	}

	return nil
}
