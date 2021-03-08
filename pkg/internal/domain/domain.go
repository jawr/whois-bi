package domain

import (
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
)

type Domain struct {
	ID int `sql:",pk"`

	Domain string `sql:",notnull,unique:domain_owner_id"`

	OwnerID int       `sql:",notnull,unique:domain_owner_id"`
	Owner   user.User `sql:"fk:owner_id" json:"-"`

	// meta data
	AddedAt   JsonDate    `sql:",type:timestamptz,notnull,default:now()"`
	DeletedAt pg.NullTime `pg:",type:timestamptz,soft_delete"`

	// when was this domain last updated, useful for starting jobs
	LastJobAt     JsonDate `sql:",type:timestamptz,null"`
	LastUpdatedAt JsonDate `sql:",type:timestamptz,null"`
}

// create a new domain attached to an owner
func NewDomain(domain string, owner user.User) Domain {
	domain = strings.TrimSpace(domain)

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
