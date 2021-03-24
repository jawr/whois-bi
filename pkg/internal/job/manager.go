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
	"github.com/go-pg/pg"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/sender"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Manager struct {
	db      *pg.DB
	router  *message.Router
	emailer *sender.Sender

	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewManager(db *pg.DB, emailer *sender.Sender) (*Manager, error) {
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
		return nil, errors.WithMessage(err, "NewPublisher")
	}

	// setup subscriber
	subscriber, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		return nil, errors.WithMessage(err, "NewSubscriber")
	}

	routerConfig := message.RouterConfig{}

	// setup router
	router, err := message.NewRouter(routerConfig, logger)
	if err != nil {
		return nil, errors.WithMessage(err, "NewRouter")
	}

	// setup manager
	manager := Manager{
		db:         db,
		emailer:    emailer,
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
		// recover from any panics
		middleware.Recoverer,
	)

	router.AddNoPublisherHandler(
		"job.response",
		"job.response",
		subscriber,
		manager.jobResponseHandler(),
	)

	return &manager, nil
}

func (m *Manager) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	// handle creation of jobs
	eg.Go(func() error {
		return m.createJobs(ctx)
	})

	// run router
	eg.Go(func() error {
		return m.router.Run(ctx)
	})

	// handle alerts
	eg.Go(func() error {
		return m.sendAlerts(ctx)
	})

	if err := eg.Wait(); err != nil {
		log.Printf("Error encountered: %s", err)
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
			return errors.WithMessage(err, "GetdomainsWhereLastJobBefore")
		}

		log.Printf("Found %d domains", len(domains))

		for _, d := range domains {
			log.Printf("CREATE JOB FOR %s", d.Domain)
			j := NewJob(d)
			if err := j.Insert(m.db); err != nil {
				continue
			}
		}

		// get all jobs without a started_at

		jobs, err := GetJobs(m.db)
		if err != nil {
			return errors.WithMessage(err, "GetJobs")
		}

		log.Printf("Found %d jobs", len(jobs))

		for _, j := range jobs {
			log.Printf("job: %+v", j)

			currentRecords, err := j.Domain.GetRecords(m.db)
			if err != nil {
				return errors.WithMessage(err, "GetRecords")
			}
			j.CurrentRecords = currentRecords

			b, err := json.Marshal(&j)
			if err != nil {
				return errors.WithMessage(err, "Marhsal")
			}

			msg := message.NewMessage(watermill.NewUUID(), b)

			if err := m.publisher.Publish("job.queue", msg); err != nil {
				return errors.WithMessage(err, "Publish")
			}
		}

		if len(jobs) > 0 {
			if _, err := m.db.Model(&jobs).Set("started_at = now()").WherePK().Update(); err != nil {
				return errors.WithMessage(err, "Update started_at")
			}
		}
	}
	return nil
}

func (m *Manager) jobResponseHandler() message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		var response JobResponse
		if err := json.Unmarshal(msg.Payload, &response); err != nil {
			return errors.WithMessage(err, "Unmarshal")
		}

		log.Printf("JOB RESPONSE %d / %s", response.Job.ID, response.Job.Domain)

		response.OwnerID = response.Job.Domain.OwnerID

		log.Printf(
			"Found %d additions and %d removals",
			len(response.RecordAdditions),
			len(response.RecordRemovals),
		)

		log.Println("Removals: ")
		for _, record := range response.RecordRemovals {
			log.Printf("\t%d\t%s", record.ID, record.Raw)
		}

		// handle removals
		if err := response.RecordRemovals.Remove(m.db); err != nil {
			return errors.WithMessage(err, "removals.Remove")
		}

		// handle additions
		if err := response.RecordAdditions.Insert(m.db); err != nil {
			return errors.WithMessage(err, "additions.Insert")
		}

		// handle whois
		if response.Whois.Raw != nil {
			if err := response.Whois.Insert(m.db); err != nil {
				// means we hit a dupe
				if err != pg.ErrNoRows {
					log.Println(errors.WithMessage(err, "inserting whois"))
				}
			} else {
				response.WhoisUpdated = true
			}
		}

		// parse the record additions and removals through our lists to avoid sending alarm bells
		if err := m.handleLists(&response); err != nil {
			log.Printf("Error parsing lists: %s", err)
		}

		// handle alert message
		if len(response.RecordAdditions) > 0 || len(response.RecordRemovals) > 0 {
			a := Alert{
				OwnerID:  response.OwnerID,
				Response: response,
			}

			if response.Job.Domain.DontBatch {
				if err := m.handleAlerts([]Alert{a}); err != nil {
					log.Printf("Error handling alerts: %s", err)
				}
			} else {
				if _, err := m.db.Model(&a).Insert(); err != nil {
					log.Printf("Error inserting alert: %s", err)
				}
			}
		}

		_, err := m.db.Model(&response.Job).
			Set(
				"errors = ?, started_at = ?, finished_at = ?, additions = ?, removals = ?, whois_updated = ?",
				response.Errors,
				response.Job.StartedAt,
				response.Job.FinishedAt,
				len(response.RecordAdditions),
				len(response.RecordRemovals),
				response.WhoisUpdated,
			).
			WherePK().
			Update()

		if err != nil {
			return errors.WithMessage(err, "Update Job")
		}

		_, err = m.db.Model(&response.Job.Domain).Set("last_updated_at = now()").WherePK().Update()
		if err != nil {
			return errors.WithMessage(err, "Update Domain")
		}

		return nil
	}
}
