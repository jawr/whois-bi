package job

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/queue"
	"github.com/jawr/whois-bi/pkg/internal/queue/rabbit"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"
)

var (
	subdomainsToCheck = []string{
		"", "www", "mx", "media", "assets", "dashboard", "api",
		"cdn", "download", "downloads", "mail", "applytics", "email", "app",
		"img", "default._domainkey",
	}
)

type Worker struct {
	publisher queue.Publisher
	consumer  queue.Consumer

	dnsClient dns.Client
}

func NewWorker(addr string) (*Worker, error) {
	publisher := rabbit.NewPublisher(addr)
	consumer := rabbit.NewConsumer("", "job.queue", addr)

	worker := Worker{
		publisher: publisher,
		consumer:  consumer,
	}

	return &worker, nil
}

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

func (w *Worker) handleJob(ctx context.Context, msg *amqp.Delivery) {
	var job Job

	if err := json.Unmarshal(msg.Body, &job); err != nil {
		return
	}

	job.StartedAt = time.Now()

	log.Printf("JOB / %d / %s", job.ID, job.Domain.Domain)

	records := w.handleRecords(&job)

	additions, removals, err := job.Domain.CheckDelta(&w.dnsClient, job.CurrentRecords, records)
	if err != nil {
		job.Errors = append(job.Errors, errors.Wrap(err, "CheckDelta").Error())
	}

	job.RecordAdditions = additions
	job.RecordRemovals = removals
	job.FinishedAt = time.Now()

	err = w.publisher.Publish(ctx, "job.response", &job)
	if err != nil {
		log.Printf("Error handling job %d, unable to publish: %s", job.ID, err)
		return
	}
}

// build a list of current records using ANY and enumeration
func (w *Worker) handleRecords(job *Job) domain.Records {
	records, err := job.Domain.QueryANY(&w.dnsClient, job.Domain.Domain)
	if err != nil {
		job.Errors = append(job.Errors, errors.Wrap(err, "QueryANY").Error())
	}

	// explain what this does, or give it a better name
	enumRecords, err := job.Domain.QueryEnumerate(&w.dnsClient, subdomainsToCheck)
	if err != nil {
		job.Errors = append(job.Errors, errors.Wrap(err, "QueryEnumerate").Error())
	}

	if len(enumRecords) > 0 {
		records = append(records, enumRecords...)
	}

	return records
}
