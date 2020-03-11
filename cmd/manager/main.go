package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jawr/monere/job"
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

	manager, err := job.NewManager(db)
	if err != nil {
		return errors.Wrap(err, "NewManager")
	}
	defer manager.Close()

	if err := manager.Run(context.TODO()); err != nil {
		return errors.Wrap(err, "Run")
	}

	return nil
}
