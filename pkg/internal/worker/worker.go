package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jawr/whois-bi/pkg/internal/dns"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/jawr/whois-bi/pkg/internal/queue"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// A Worker consumes attempts to find record additions and removals
// for domains that are pushed on to its queue, it then publishes the
// results
type Worker struct {
	publisher queue.Publisher
	consumer  queue.Consumer

	dnsClient dns.Client
}

// NewWorker creates a worker using the provided dnsClient, publisher
// and consumer
func NewWorker(dnsClient dns.Client, publisher queue.Publisher, consumer queue.Consumer) *Worker {
	return &Worker{
		dnsClient: dnsClient,
		publisher: publisher,
		consumer:  consumer,
	}
}

// Run starts the consumer and publishers stopping on any encountered
// errors, or when the context is cancelled
func (w *Worker) Run(ctx context.Context) error {
	var wg errgroup.Group

	wg.Go(func() error {
		return w.publisher.Run(ctx)
	})

	wg.Go(func() error {
		return w.consumer.Run(ctx, w.handleJob)
	})

	return wg.Wait()
}

// handleJob decodes a Job and attempts to process it. Any errors
// encountered are pushed on to the response's Errors field
func (w *Worker) handleJob(ctx context.Context, body []byte) {
	var job job.Job

	if err := json.Unmarshal(body, &job); err != nil {
		return
	}

	job.StartedAt = time.Now()

	live, err := w.dnsClient.GetLive(
		job.Domain,
		job.CurrentRecords,
	)
	if err != nil {
		job.Errors = append(
			job.Errors,
			errors.Wrap(err, "GetLive").Error(),
		)
	} else {
		// only proceed if we had no errors otherwise we will
		// remove everything
		additions, removals := delta(job.CurrentRecords, live)

		job.RecordAdditions = additions
		job.RecordRemovals = removals
	}

	job.FinishedAt = time.Now()

	err = w.publisher.Publish(ctx, "job.response", &job)
	if err != nil {
		log.Printf("Error handling job %d, unable to publish: %s", job.ID, err)
		return
	}
}
