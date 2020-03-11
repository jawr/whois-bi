package main

import (
	"fmt"
	"os"

	"github.com/jawr/monere/api"
	"github.com/jawr/monere/sender"
	"github.com/jawr/monere/shared"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := shared.SetupDatabase()
	if err != nil {
		return errors.Wrap(err, "setupDatabase")
	}

	emailer := sender.NewSender()

	server := api.NewServer(db, emailer)

	if err := server.Run("localhost:8444"); err != nil {
		return errors.Wrap(err, "Run")
	}

	return nil
}
