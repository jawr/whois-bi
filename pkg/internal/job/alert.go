package job

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
)

const (
	// minimum amount of domains to send to 5
	alertMinDomains = 5
	// maximum age before minimum amount is irrelevant
	alertMaxAge time.Duration = time.Minute * 60
)

type Alert struct {
	ID int `pg:",pk"`

	OwnerID int       `pg:",notnull"`
	Owner   user.User `pg:"fk:owner_id,rel:has-one"`

	Response Job

	CreatedAt time.Time `pg:",notnull,default:now()"`
}

type ExpirationAlert struct {
	ID int `pg:",pk"`

	DomainID int           `pg:",notnull"`
	Domain   domain.Domain `pg:"fk:domain_id,rel:has-one"`
	SentAt   time.Time     `pg:",notnull,default:now()"`
}

func (m *Manager) sendAlerts(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// get all alerts
			var alerts []Alert
			if err := m.db.Model(&alerts).Select(); err != nil {
				return err
			}

			// sort in to owner batches
			sorted := make(map[int][]Alert, 0)

			for _, a := range alerts {
				if _, ok := sorted[a.OwnerID]; !ok {
					sorted[a.OwnerID] = []Alert{}
				}
				sorted[a.OwnerID] = append(sorted[a.OwnerID], a)
			}

			for owner, alerts := range sorted {
				// check minimum threshold
				if len(alerts) < alertMinDomains {
					var handle bool
					for _, a := range alerts {
						if time.Since(a.CreatedAt) > alertMaxAge {
							handle = true
							break
						}
					}

					if !handle {
						continue
					}
				}

				if err := m.handleAlerts(alerts); err != nil {
					log.Printf("Error handling alerts for owner %d: %s", owner, err)
				}

				if _, err := m.db.Model(&alerts).Delete(); err != nil {
					return err
				}
			}

			// find domains about to expire
			var whois []domain.Whois

			err := m.db.Model(&whois).
				DistinctOn("whois.domain_id").
				Join("LEFT JOIN expiration_alerts AS ea ON ea.domain_id = whois.domain_id AND date_trunc('day', ea.sent_at) = date_trunc('day', NOW())").
				Relation("Domain").
				Where(
					`
						date_trunc('day', whois.expiration_date) = date_trunc('day', ?::timestamp) 
						AND ea.id IS NULL
					`,
					time.Now().AddDate(0, 0, 7),
				).Select()
			if err != nil {
				return err
			}

			if err := m.handleExpirationAlerts(whois); err != nil {
				log.Printf("Error handling expiration alerts: %s", err)
			}

			if len(whois) > 0 {
				log.Printf("Whois alert for %d domains", len(whois))
			}
		}
	}

	return nil
}

func (m *Manager) handleExpirationAlerts(whois []domain.Whois) error {
	for _, w := range whois {
		var ea = ExpirationAlert{
			DomainID: w.DomainID,
		}
		if _, err := m.db.Model(&ea).Insert(); err != nil {
			return err
		}

		subject := fmt.Sprintf("ALARM BELLS - %s expires in 7 days", w.Domain.Domain)

		body := fmt.Sprintf(
			"Your domain will expire in 7 days for more information visit: https://%s/domain/%s\n\n",
			os.Getenv("DOMAIN"),
			w.Domain.Domain,
		)

		var owner user.User

		if err := m.db.Model(&owner).Where("id = ?", w.Domain.OwnerID).Select(); err != nil {
			return errors.WithMessage(err, "Select Owner")
		}

		if err := m.emailer.Send(owner.Email, subject, body); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) handleAlerts(alerts []Alert) error {
	subject := fmt.Sprintf("ALARM BELLS - Changes to %d domains", len(alerts))

	var body strings.Builder

	var ownerID int

	for _, alert := range alerts {
		response := alert.Response

		fmt.Fprintf(&body, "<pre>")

		fmt.Fprintf(
			&body,
			"New changes have been detected, please go to: https://%s/domain/%s for more details or find a summary of the changes below.\n\n",
			os.Getenv("DOMAIN"),
			response.Domain.Domain,
		)

		if response.WhoisUpdated {
			fmt.Fprintf(
				&body,
				"Whois has been updated!\n\n",
			)
		}

		for idx, record := range response.RecordAdditions {
			if idx == 0 {
				fmt.Fprintf(&body, "-------------------------------- / additions start\n")
			}
			fmt.Fprintf(&body, "\t+++\t%s\n", record.Raw)
		}

		for idx, record := range response.RecordRemovals {
			if idx == 0 {
				fmt.Fprintf(&body, "-------------------------------- / removals start\n")
			}
			fmt.Fprintf(&body, "\t---\t%s\n", record.Raw)
		}

		fmt.Fprintf(&body, "-------------------------------- / end\n")

		fmt.Fprintf(&body, "</pre>")

		if ownerID == 0 {
			ownerID = response.Domain.OwnerID
		}
	}

	if ownerID == 0 {
		return nil
	}

	var owner user.User

	if err := m.db.Model(&owner).Where("id = ?", ownerID).Select(); err != nil {
		return errors.WithMessage(err, "Select Owner")
	}

	if err := m.emailer.Send(owner.Email, subject, body.String()); err != nil {
		return err
	}

	return nil
}
