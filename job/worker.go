package job

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/jawr/monere/domain"
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

	amqpConfig := amqp.NewDurableQueueConfig(
		"amqp://172.17.0.2:5672/",
	)

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

	// setup router
	router, err := message.NewRouter(message.RouterConfig{}, logger)
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
		"monere.job.queue",
		"monere.job.queue",
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
	client := dns.Client{
		Net: "tcp",
	}

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

			if err := w.publisher.Publish("monere.job.response", msg); err != nil {
				finalErr = errors.Wrap(err, "Publish")
				return
			}

			if len(response.Error) > 0 {
				finalErr = errors.New(response.Error)
				return
			}

		}()

		log.Printf("JOB / %d", job.DomainID)

		records, err := job.Domain.QueryANY(&client, job.Domain.Domain)
		if err != nil {
			response.Error = errors.Wrap(err, "QueryANY").Error()
			return
		}

		additions, removals, err := job.Domain.CheckDelta(&client, records)
		if err != nil {
			response.Error = errors.Wrap(err, "CheckDelta").Error()
			return
		}

		response.RecordAdditions = additions
		response.RecordRemovals = removals

		w, err := domain.NewWhois(job.Domain)
		if err != nil {
			response.Error = errors.Wrap(err, "NewWhois").Error()
			return
		}

		response.Whois = w

		return nil
	}
}
