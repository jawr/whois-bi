package job

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/pkg/errors"
)

type Job struct {
	ID int `sql:",pk" json:"id"`

	DomainID int           `sql:",notnull" json:"domain_id"`
	Domain   domain.Domain `sql:"fk:domain_id"`

	Errors []string `sql:",notnull" json:"errors"`

	Additions    int  `sql:",notnull" json:"additions"`
	Removals     int  `sql:",notnull" json:"removals"`
	WhoisUpdated bool `sql:",notnull" json:"whois_updated"`

	CreatedAt  time.Time `sql:",notnull,default:now()" json:"created_at"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`

	// dont persist, what is this for?
	CurrentRecords domain.Records `sql:"-" json:"current_records"`
}

type JobResponse struct {
	Job Job

	Errors []string

	RecordAdditions domain.Records
	RecordRemovals  domain.Records
	Whois           domain.Whois
}

func NewJob(d domain.Domain) Job {
	j := Job{
		DomainID: d.ID,
		Domain:   d,
	}

	return j
}

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

// find all jobs that have not yet started
func GetJobs(db *pg.DB) ([]Job, error) {
	var jobs []Job
	err := db.Model(&jobs).Relation("Domain").Where("started_at IS NULL OR (started_at < NOW() - INTERVAL '1 HOUR' AND finished_at IS NULL)").Select()
	if err != nil {
		return nil, err
	}
	return jobs, nil
}
