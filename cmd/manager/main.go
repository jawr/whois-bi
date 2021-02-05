package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jawr/whois.bi/internal/cmdutil"
	"github.com/jawr/whois.bi/internal/job"
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
