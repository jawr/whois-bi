package main

import (
	"fmt"
	"os"

	"github.com/jawr/whois-bi/pkg/internal/api"
	"github.com/jawr/whois-bi/pkg/internal/db"
	"github.com/jawr/whois-bi/pkg/internal/emailer"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dbConn, err := db.SetupDatabase()
	if err != nil {
		return errors.WithMessage(err, "SetupDatabase")
	}
	defer dbConn.Close()

	sender := emailer.NewSMTPSenderFromEnv()
	emailer, err := emailer.NewEmailer(
		os.Getenv("SMTP_FROM_NAME"),
		os.Getenv("SMTP_EMAIL"),
		sender,
	)
	if err != nil {
		return errors.WithMessage(err, "NewEmailer")
	}

	server := api.NewServer(dbConn, emailer)

	if err := server.Run(os.Getenv("HTTP_API_ADDR")); err != nil {
		return errors.Wrap(err, "Run")
	}

	return nil
}
