package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jawr/monere/job"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	worker, err := job.NewWorker()
	if err != nil {
		return errors.Wrap(err, "NewWorker")
	}
	defer worker.Close()

	if err := worker.Run(context.TODO()); err != nil {
		return errors.Wrap(err, "Run")
	}

	return nil
}
