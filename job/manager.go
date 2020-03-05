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
	"github.com/go-pg/pg"
	"github.com/jawr/monere/domain"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Manager struct {
	db     *pg.DB
	router *message.Router

	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewManager(db *pg.DB) (*Manager, error) {
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

	routerConfig := message.RouterConfig{
		CloseTimeout: time.Second * 30,
	}

	// setup router
	router, err := message.NewRouter(routerConfig, logger)
	if err != nil {
		return nil, errors.Wrap(err, "NewRouter")
	}

	// setup manager
	manager := Manager{
		db:         db,
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
		"monere.job.response",
		"monere.job.response",
		subscriber,
		manager.jobResponseHandler(),
	)

	return &manager, nil
}

func (m *Manager) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	eg, ctx := errgroup.WithContext(ctx)

	// handle creation of jobs
	eg.Go(func() error {
		defer cancel()
		return m.createJobs(ctx)
	})

	// run router
	eg.Go(func() error {
		defer cancel()
		return m.router.Run(ctx)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Close() error {
	return m.router.Close()
}

func (m *Manager) createJobs(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}

		// search for new jobs and dispatch them
		domains, err := domain.GetDomainsWhereLastJobBefore(m.db, time.Hour*24)
		if err != nil {
			return errors.Wrap(err, "GetdomainsWhereLastJobBefore")
		}

		for _, d := range domains {
			log.Printf("CREATE JOB FOR %s", d.Domain)
			j := NewJob(d)
			if err := j.Insert(m.db); err != nil {
				return errors.Wrap(err, "Insert job")
			}

			currentRecords, err := j.Domain.GetRecords(m.db)
			if err != nil {
				return errors.Wrap(err, "GetRecords")
			}
			j.CurrentRecords = currentRecords

			b, err := json.Marshal(&j)
			if err != nil {
				return errors.Wrap(err, "Marhsal")
			}

			msg := message.NewMessage(watermill.NewUUID(), b)

			if err := m.publisher.Publish("monere.job.queue", msg); err != nil {
				return errors.Wrap(err, "Publish")
			}

		}
	}
	return nil
}

func (m *Manager) jobResponseHandler() message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		var response JobResponse
		if err := json.Unmarshal(msg.Payload, &response); err != nil {
			return errors.Wrap(err, "Unmarshal")
		}

		log.Printf("JOB RESPONSE %d", response.Job.ID)

		if len(response.Error) > 0 {
			log.Printf("Error: %s", response.Error)

			_, err := m.db.Model(&response.Job).
				Set(
					"started_at = ? AND finished_at = ? AND error = ?",
					response.Job.StartedAt,
					response.Job.FinishedAt,
					response.Error,
				).
				WherePK().
				Update()

			if err != nil {
				return errors.Wrap(err, "Update Job Error")
			}
			return nil
		}

		log.Printf(
			"Found %d additions and %d removals",
			len(response.RecordAdditions),
			len(response.RecordRemovals),
		)

		// handle removals
		if err := response.RecordRemovals.Remove(m.db); err != nil {
			return errors.Wrap(err, "removals.Remove")
		}

		// handle additions
		if err := response.RecordAdditions.Insert(m.db); err != nil {
			return errors.Wrap(err, "additions.Insert")
		}

		for _, record := range response.RecordAdditions {
			log.Printf("\t+++\t%s", record)
		}

		for _, record := range response.RecordRemovals {
			log.Printf("\t---\t%s", record)
		}

		// handle whois
		whoisUpdated := true
		if err := response.Whois.Insert(m.db); err != nil {
			whoisUpdated = false
			// return errors.Wrap(err, "Insert whois")
		}

		_, err := m.db.Model(&response.Job).
			Set(
				"started_at = ?, finished_at = ?, additions = ?, removals = ?, whois_updated = ?",
				response.Job.StartedAt,
				response.Job.FinishedAt,
				len(response.RecordAdditions),
				len(response.RecordRemovals),
				whoisUpdated,
			).
			WherePK().
			Update()

		if err != nil {
			return errors.Wrap(err, "Update Job")
		}

		_, err = m.db.Model(&response.Job.Domain).Set("last_updated_at = now()").WherePK().Update()
		if err != nil {
			return errors.Wrap(err, "Update Domain")
		}

		return nil
	}
}
