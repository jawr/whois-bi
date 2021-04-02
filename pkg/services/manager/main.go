package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jawr/whois-bi/pkg/internal/db"
	"github.com/jawr/whois-bi/pkg/internal/emailer"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/jawr/whois-bi/pkg/internal/queue/rabbit"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	addr := os.Getenv("RABBITMQ_URI")
	if len(addr) == 0 {
		return errors.New("No RABBITMQ_URI")
	}

	publisher := rabbit.NewPublisher(addr)
	consumer := rabbit.NewConsumer("", "job.response", addr)

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

	manager, err := job.NewManager(publisher, consumer, dbConn, emailer)
	if err != nil {
		return errors.WithMessage(err, "NewManager")
	}

	ctx := context.Background()

	if err := manager.Run(ctx); err != nil {
		return errors.WithMessage(err, "Run")
	}

	return nil
}
