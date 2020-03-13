package domain

import (
	"encoding/json"
	"time"

	"github.com/go-pg/pg"
	"github.com/jawr/monere/user"
	"github.com/pkg/errors"
)

type Domain struct {
	ID int `sql:",pk"`

	Domain string `sql:",notnull,unique:domain_owner_id"`

	OwnerID int       `sql:",notnull,unique:domain_owner_id"`
	Owner   user.User `sql:"fk:owner_id" json:"-"`

	// meta data
	AddedAt   time.Time `sql:",notnull,default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`

	// when was this domain last updated, useful for starting jobs
	LastJobAt     time.Time
	LastUpdatedAt time.Time
}

// create a new domain attached to an owner
func NewDomain(domain string, owner user.User) Domain {
	newDomain := Domain{
		Domain:  domain,
		OwnerID: owner.ID,
		Owner:   owner,
	}

	return newDomain
}

// get domain by name
func GetDomain(db *pg.DB, domain string) (Domain, error) {
	var dom Domain
	if err := db.Model(&dom).Where("domain = ?", domain).Select(); err != nil {
		return Domain{}, err
	}
	return dom, nil
}

// get domains where lastUpdatedAt > d
func GetDomainsWhereLastJobBefore(db *pg.DB, d time.Duration) ([]Domain, error) {
	var domains []Domain
	before := time.Now().Add(d * -1)
	if err := db.Model(&domains).Where("last_job_at IS NULL OR last_job_at < ?", before).Select(); err != nil {
		return nil, err
	}
	return domains, nil
}

// insert a domain in to the database
func (d *Domain) Insert(db *pg.DB) error {
	if _, err := db.Model(d).Returning("*").Insert(); err != nil {
		return err
	}

	return nil
}

// get most recent records for a domain
func (d Domain) GetRecords(db *pg.DB) (Records, error) {
	var records Records
	err := db.Model(&records).
		Where(
			"domain_id = ? AND removed_at IS NULL",
			d.ID,
		).
		Select()
	if err != nil {
		return nil, errors.Wrap(err, "Select records")
	}

	return records, nil
}

func (d Domain) String() string {
	return d.Domain
}

// custom marshaller
func (d *Domain) MarshalJSON() ([]byte, error) {
	type Alias Domain

	// can be null
	var deletedAt, lastJobAt, lastUpdatedAt string
	if !d.DeletedAt.IsZero() {
		deletedAt = d.DeletedAt.Format("2006/01/02 15:04")
	}
	if !d.LastJobAt.IsZero() {
		lastJobAt = d.LastJobAt.Format("2006/01/02 15:04")
	}
	if !d.LastUpdatedAt.IsZero() {
		lastUpdatedAt = d.LastUpdatedAt.Format("2006/01/02 15:04")
	}

	return json.Marshal(&struct {
		AddedAt       string
		DeletedAt     string
		LastJobAt     string
		LastUpdatedAt string
		*Alias
	}{
		AddedAt:       d.AddedAt.Format("2006/01/02 15:04"),
		DeletedAt:     deletedAt,
		LastJobAt:     lastJobAt,
		LastUpdatedAt: lastUpdatedAt,
		Alias:         (*Alias)(d),
	})
}

func (d *Domain) UnmarshalJSON(data []byte) error {
	type Alias Domain

	aux := &struct {
		AddedAt       string
		DeletedAt     string
		LastJobAt     string
		LastUpdatedAt string
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error

	d.AddedAt, err = time.Parse("2006/01/02 15:04", aux.AddedAt)
	if err != nil {
		return err
	}
	if len(aux.DeletedAt) > 0 {
		d.DeletedAt, err = time.Parse("2006/01/02 15:04", aux.DeletedAt)
		if err != nil {
			return err
		}
	}
	if len(aux.LastJobAt) > 0 {
		d.LastJobAt, err = time.Parse("2006/01/02 15:04", aux.LastJobAt)
		if err != nil {
			return err
		}
	}
	if len(aux.LastUpdatedAt) > 0 {
		d.LastUpdatedAt, err = time.Parse("2006/01/02 15:04", aux.LastUpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil

}
