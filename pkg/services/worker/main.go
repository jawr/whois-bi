package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

	worker, err := job.NewWorker(addr)
	if err != nil {
		return errors.WithMessage(err, "NewWorker")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg errgroup.Group

	wg.Go(func() error {
		return worker.Run(ctx)
	})

	wg.Go(func() error {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)

		select {
		case <-sigc:
			cancel()
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	return wg.Wait()
}
