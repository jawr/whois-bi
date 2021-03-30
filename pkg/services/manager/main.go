package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jawr/whois-bi/pkg/internal/cmdutil"
	"github.com/jawr/whois-bi/pkg/internal/emailer"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := cmdutil.SetupDatabase()
	if err != nil {
		return errors.WithMessage(err, "SetupDatabase")
	}
	defer db.Close()

	sender := emailer.NewSMTPSenderFromEnv()
	emailer, err := emailer.NewEmailer(
		os.Getenv("SMTP_FROM_NAME"),
		os.Getenv("SMTP_EMAIL"),
		sender,
	)
	if err != nil {
		return errors.WithMessage(err, "NewEmailer")
	}

	manager, err := job.NewManager(db, emailer)
	if err != nil {
		return errors.WithMessage(err, "NewManager")
	}
	defer manager.Close()

	ctx := context.Background()

	if err := manager.Run(ctx); err != nil {
		return errors.WithMessage(err, "Run")
	}

	return nil
}
