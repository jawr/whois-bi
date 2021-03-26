package job

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/pkg/errors"
)

type Job struct {
	ID int `pg:",pk" json:"id"`

	DomainID int           `pg:",notnull" json:"domain_id"`
	Domain   domain.Domain `pg:"fk:domain_id,rel:has-one"`

	Errors []string `pg:",use_zero" json:"errors"`

	Additions    int  `pg:",use_zero" json:"additions"`
	Removals     int  `pg:",use_zero" json:"removals"`
	WhoisUpdated bool `pg:",use_zero" json:"whois_updated"`

	CreatedAt  time.Time `pg:",notnull,default:now()" json:"created_at"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`

	// dont persist, what is this for?
	CurrentRecords domain.Records `pg:"-" json:"current_records"`
}

type JobResponse struct {
	OwnerID int

	Job Job

	Errors []string

	RecordAdditions domain.Records
	RecordRemovals  domain.Records
	Whois           domain.Whois

	WhoisUpdated bool
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
