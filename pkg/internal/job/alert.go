package job

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	ID int `sql:",pk"`

	OwnerID int       `sql:",notnull"`
	Owner   user.User `sql:"fk:owner_id"`

	Response JobResponse

	CreatedAt time.Time `sql:",notnull,default:now()"`
}

func (m *Manager) sendAlerts(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:

		default:
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

		}
	}

	return nil
}

func (m *Manager) handleAlerts(alerts []Alert) error {
	alertSubject := fmt.Sprintf("ALARM BELLS - Changes to %d domains", len(alerts))

	var alertBody strings.Builder

	var ownerID int

	for _, alert := range alerts {
		response := alert.Response

		fmt.Fprintf(&alertBody, "<pre>")

		fmt.Fprintf(
			&alertBody,
			"New changes have been detected, please go to: https://%s/#/dashboard/%s for more details or find a summary of the changes below.\n\n",
			os.Getenv("DOMAIN"),
			response.Job.Domain.Domain,
		)

		if response.WhoisUpdated {
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

		if ownerID == 0 {
			ownerID = response.OwnerID
		}
	}

	if ownerID == 0 {
		return nil
	}

	var owner user.User

	if err := m.db.Model(&owner).Where("id = ?", ownerID).Select(); err != nil {
		return errors.WithMessage(err, "Select Owner")
	}

	if err := m.emailer.Send(owner.Email, alertSubject, alertBody.String()); err != nil {
		return err
	}

	return nil
}
