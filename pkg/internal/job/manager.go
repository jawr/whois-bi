package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/go-pg/pg"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/sender"
	"github.com/jawr/whois-bi/pkg/internal/user"
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
		var whoisUpdated bool
		if response.Whois.Raw != nil {
			if err := response.Whois.Insert(m.db); err != nil {
				// means we hit a dupe
				if err != pg.ErrNoRows {
					log.Println(errors.WithMessage(err, "inserting whois"))
				}
			} else {
				whoisUpdated = true
			}
		}

		// handle alert message
		if len(response.RecordAdditions) > 0 || len(response.RecordRemovals) > 0 {
			alertSubject := fmt.Sprintf("ALARM BELLS - Changes to domain '%s'", response.Job.Domain.Domain)
			var alertBody strings.Builder

			fmt.Fprintf(&alertBody, "<pre>")

			fmt.Fprintf(
				&alertBody,
				"New changes have been detected, please go to: https://%s/#/dashboard/%s for more details or find a summary of the changes below.\n\n",
				os.Getenv("DOMAIN"),
				response.Job.Domain.Domain,
			)

			if whoisUpdated {
				fmt.Fprintf(
					&alertBody,
					"Whois has been updated!\n\n",
				)
			}

			for idx, record := range response.RecordAdditions {
				if idx == 0 {
					fmt.Fprintf(&alertBody, "-------------------------------- / additions start\n")
				}
				fmt.Fprintf(&alertBody, "\t+++\t%s\n", record.Raw)
			}

			for idx, record := range response.RecordRemovals {
				if idx == 0 {
					fmt.Fprintf(&alertBody, "-------------------------------- / removals start\n")
				}
				fmt.Fprintf(&alertBody, "\t---\t%s\n", record.Raw)
			}

			fmt.Fprintf(&alertBody, "-------------------------------- / end\n")

			fmt.Fprintf(&alertBody, "</pre>")

			// get user
			var owner user.User

			if err := m.db.Model(&owner).Where("id = ?", response.Job.Domain.OwnerID).Select(); err != nil {
				return errors.WithMessage(err, "Select Owner")
			}

			if err := m.emailer.Send(owner.Email, alertSubject, alertBody.String()); err != nil {
				log.Printf("Error Sending to %s: %s", owner.Email, err)
			}
		}

		_, err := m.db.Model(&response.Job).
			Set(
				"errors = ?, started_at = ?, finished_at = ?, additions = ?, removals = ?, whois_updated = ?",
				response.Job.Errors,
				response.Job.StartedAt,
				response.Job.FinishedAt,
				len(response.RecordAdditions),
				len(response.RecordRemovals),
				whoisUpdated,
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
