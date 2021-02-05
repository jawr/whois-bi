package main

import (
	"fmt"
	"os"

	"github.com/jawr/whois.bi/internal/api"
	"github.com/jawr/whois.bi/internal/cmdutil"
	"github.com/jawr/whois.bi/internal/sender"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := cmdutil.LoadDotEnv(); err != nil {
		return errors.WithMessage(err, "LoadDotEnv")
	}

	db, err := cmdutil.SetupDatabase()
	if err != nil {
		return errors.WithMessage(err, "SetupDatabase")
	}
	defer db.Close()

	emailer := sender.NewSender()

	server := api.NewServer(db, emailer)

	if err := server.Run(os.Getenv("HTTP_API_ADDR")); err != nil {
		return errors.Wrap(err, "Run")
	}

	return nil
}
