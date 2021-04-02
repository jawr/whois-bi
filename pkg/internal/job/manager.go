package job

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/emailer"
	"github.com/jawr/whois-bi/pkg/internal/queue"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Manager struct {
	db      *pg.DB
	emailer *emailer.Emailer

	publisher queue.Publisher
	consumer  queue.Consumer
}

func NewManager(publisher queue.Publisher, consumer queue.Consumer, db *pg.DB, emailer *emailer.Emailer) (*Manager, error) {

	// setup manager
	manager := Manager{
		db:        db,
		emailer:   emailer,
		publisher: publisher,
		consumer:  consumer,
	}

	return &manager, nil
}

func (m *Manager) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	// handle creation of jobs
	wg.Go(func() error {
		return m.createJobs(ctx)
	})

	wg.Go(func() error {
		return m.publisher.Run(ctx)
	})

	wg.Go(func() error {
		return m.consumer.Run(ctx, m.handleJobResponses)
	})

	// handle alerts
	wg.Go(func() error {
		return m.sendAlerts(ctx)
	})

	return wg.Wait()
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

		for _, d := range domains {
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

		if len(jobs) > 0 {
			log.Printf("Found %d jobs", len(jobs))
		}

		for _, j := range jobs {

			currentRecords, err := j.Domain.GetRecords(m.db)
			if err != nil {
				return errors.WithMessage(err, "GetRecords")
			}
			j.CurrentRecords = currentRecords

			if err := m.publisher.Publish(ctx, "job.queue", &j); err != nil {
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

func (m *Manager) handleJobResponses(ctx context.Context, body []byte) {
	var job Job
	if err := json.Unmarshal(body, &job); err != nil {
		return
	}

	log.Printf(
		"Job %d  / %s Found %d additions and %d removals",
		job.ID,
		job.Domain,
		len(job.RecordAdditions),
		len(job.RecordRemovals),
	)

	// handle removals
	if err := job.RecordRemovals.Remove(m.db); err != nil {
		log.Printf("Error RecordRemovals.Remove() job %d: %s", job.ID, err)
		return
	}

	// handle additions
	if err := job.RecordAdditions.Insert(m.db); err != nil {
		log.Printf("Error RecordAdditions.Insert() job %d: %s", job.ID, err)
		return
	}

	// handle whois
	if job.Whois.Raw != nil {
		if err := job.Whois.Insert(m.db); err != nil {
			// means we hit a dupe
			if err != pg.ErrNoRows {
				log.Println(errors.WithMessage(err, "inserting whois"))
			}
		} else {
			job.WhoisUpdated = true
		}
	}

	// parse the record additions and removals through our lists to avoid sending alarm bells
	if err := m.handleLists(&job); err != nil {
		log.Printf("Error parsing lists for job %d: %s", job.ID, err)
	}

	// handle alert message
	if len(job.RecordAdditions) > 0 || len(job.RecordRemovals) > 0 {
		a := Alert{
			OwnerID:  job.Domain.OwnerID,
			Response: job,
		}

		if job.Domain.DontBatch {
			if err := m.handleAlerts([]Alert{a}); err != nil {
				log.Printf("Error handling alerts for job %d: %s", job.ID, err)
			}
		} else {
			if _, err := m.db.Model(&a).Insert(); err != nil {
				log.Printf("Error inserting alert for job %d: %s", job.ID, err)
			}
		}
	}

	_, err := m.db.Model(&job).
		Set(
			"errors = ?, started_at = ?, finished_at = ?, additions = ?, removals = ?, whois_updated = ?",
			job.Errors,
			job.StartedAt,
			job.FinishedAt,
			len(job.RecordAdditions),
			len(job.RecordRemovals),
			job.WhoisUpdated,
		).
		WherePK().
		Update()

	if err != nil {
		log.Printf("Error updating job %d: %s", job.ID, err)
		return
	}

	_, err = m.db.Model(&job.Domain).Set("last_updated_at = now()").WherePK().Update()
	if err != nil {
		log.Printf("Error updating domain %d: %s", job.ID, err)
		return
	}
}
