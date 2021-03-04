package main

import (
	"context"
	"fmt"
	"os"

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
	worker, err := job.NewWorker()
	if err != nil {
		return errors.WithMessage(err, "NewWorker")
	}
	defer worker.Close()

	ctx := context.Background()

	if err := worker.Run(ctx); err != nil {
		return errors.WithMessage(err, "Run")
	}

	return nil
}
