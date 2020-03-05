package job

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/jawr/monere/domain"
	"github.com/pkg/errors"
)

type Job struct {
	ID int `sql:",pk"`

	DomainID int           `sql:",notnull"`
	Domain   domain.Domain `sql:"fk:domain_id"`

	Error string

	Additions    int  `sql:",notnull"`
	Removals     int  `sql:",notnull"`
	WhoisUpdated bool `sql:",notnull"`

	CreatedAt  time.Time `sql:",notnull,default:now()"`
	StartedAt  time.Time
	FinishedAt time.Time

	// dont persist
	CurrentRecords domain.Records `sql:"-"`
}

type JobResponse struct {
	Job Job

	Error string

	RecordAdditions domain.Records
	RecordRemovals  domain.Records
	Whois           domain.Whois
}

// create a new job from a domain
func NewJob(d domain.Domain) Job {
	j := Job{
		DomainID: d.ID,
		Domain:   d,
	}

	return j
}

// insert new job
func (j *Job) Insert(db *pg.DB) error {
	_, err := db.Model(j).Returning("*").Insert()
	if err != nil {
		return errors.Wrap(err, "Insert Job")
	}

	// update domain
	_, err = db.Model(&j.Domain).
		Set("last_job_at = now()").
		Where("id = ?", j.DomainID).
		Update()
	if err != nil {
		return errors.Wrap(err, "Update Domain")
	}

	return nil
}
