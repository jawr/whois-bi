package job

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type Worker struct {
	router *message.Router

	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewWorker() (*Worker, error) {
	logger := watermill.NewStdLogger(false, false)

	addr := os.Getenv("RABBITMQ_URI")
	if len(addr) == 0 {
		return nil, errors.New("No RABBITMQ_URI")
	}

	amqpConfig := amqp.NewDurableQueueConfig(addr)
	amqpConfig.Consume.NoRequeueOnNack = true

	// setup publisher
	publisher, err := amqp.NewPublisher(amqpConfig, logger)
	if err != nil {
		return nil, errors.Wrap(err, "NewPublisher")
	}

	// setup subscriber
	subscriber, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		return nil, errors.Wrap(err, "NewSubscriber")
	}

	routerConfig := message.RouterConfig{
		CloseTimeout: time.Second * 30,
	}

	// setup router
	router, err := message.NewRouter(routerConfig, logger)
	if err != nil {
		return nil, errors.Wrap(err, "NewRouter")
	}

	worker := Worker{
		router:     router,
		publisher:  publisher,
		subscriber: subscriber,
	}

	// finish setting up router
	router.AddPlugin(
		// gracefully handle SIGTERM
		plugin.SignalsHandler,
	)

	router.AddMiddleware(
		// forward correlation ids to produced messages
		middleware.CorrelationID,
		// recover from any panics
		middleware.Recoverer,
	)

	router.AddNoPublisherHandler(
		"job.queue",
		"job.queue",
		subscriber,
		worker.jobHandler(),
	)

	return &worker, nil
}

func (w *Worker) Run(ctx context.Context) error {
	return w.router.Run(ctx)
}

func (w *Worker) Close() error {
	return w.router.Close()
}

func (w *Worker) jobHandler() message.NoPublishHandlerFunc {
	// make further upstream?
	client := dns.Client{}

	return func(msg *message.Message) (finalErr error) {
		var job Job
		if err := json.Unmarshal(msg.Payload, &job); err != nil {
			return errors.Wrap(err, "Unmarshal")
		}

		job.StartedAt = time.Now()
		response := JobResponse{
			Job: job,
		}

		// always dispatch our response
		defer func() {
			response.Job.FinishedAt = time.Now()

			b, err := json.Marshal(&response)
			if err != nil {
				finalErr = errors.Wrap(err, "Marshal")
				return
			}

			msg := message.NewMessage(watermill.NewUUID(), b)

			if err := w.publisher.Publish("job.response", msg); err != nil {
				finalErr = errors.Wrap(err, "Publish")
				return
			}

			if len(response.Errors) > 0 {
				finalErr = errors.Errorf("%d errors encountered", len(response.Errors))
			}

		}()

		log.Printf("JOB / %d / %s", job.ID, job.Domain.Domain)

		records, err := job.Domain.QueryANY(&client, job.Domain.Domain)
		if err != nil {
			response.Errors = append(response.Errors, errors.Wrap(err, "QueryANY").Error())
		}

		enumRecords, err := job.Domain.QueryEnumerate(&client, []string{
			"", "www", "mx", "media", "assets", "dashboard", "api",
			"cdn", "download", "downloads", "mail", "applytics", "email", "app",
			"img", "default._domainkey",
		})
		if err != nil {
			response.Errors = append(response.Errors, errors.Wrap(err, "QueryEnumerate").Error())
		}

		log.Printf("Found %d enumRecords", len(enumRecords))

		records = append(records, enumRecords...)

		additions, removals, err := job.Domain.CheckDelta(&client, job.CurrentRecords, records)
		if err != nil {
			response.Errors = append(response.Errors, errors.Wrap(err, "CheckDelta").Error())
		}

		response.RecordAdditions = additions
		response.RecordRemovals = removals

		w, err := domain.NewWhois(job.Domain)
		if err != nil {
			response.Errors = append(response.Errors, errors.Wrap(err, "NewWhois").Error())
		} else {
			response.Whois = w
		}

		return nil
	}
}
