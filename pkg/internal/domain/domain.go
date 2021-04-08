package domain

import (
	"strings"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/go-pg/pg/v10"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
)

type Domain struct {
	ID int `pg:",pk" json:"id"`

	Domain string `pg:",notnull,unique:domain_owner_id" json:"domain"`

	OwnerID int       `pg:",notnull,unique:domain_owner_id" json:"owner_id"`
	Owner   user.User `pg:"fk:owner_id,rel:has-one" json:"-"`

	// settings
	DontBatch bool `pg:",notnull,use_zero" json:"dont_batch"`

	// meta data
	AddedAt   time.Time   `pg:",type:timestamptz,notnull,default:now()" json:"added_at"`
	DeletedAt pg.NullTime `pg:",type:timestamptz,soft_delete" json:"deleted_at"`

	// when was this domain last updated, useful for starting jobs
	LastJobAt     time.Time `pg:",type:timestamptz" json:"last_job_at"`
	LastUpdatedAt time.Time `pg:",type:timestamptz" json:"last_updated_at"`
}

// create a new domain attached to an owner
func NewDomain(domain string, owner user.User) Domain {
	domain = strings.TrimSpace(domain)

	// validation
	domain = domainutil.Domain(domain)

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
