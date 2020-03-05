package domain

import (
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/jawr/monere/user"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type Domain struct {
	ID int `sql:",pk"`

	Domain string `sql:",notnull,unique"`

	OwnerID int       `sql:",notnull"`
	Owner   user.User `sql:"fk:owner_id"`

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

// perform an any query
func (d Domain) QueryANY(client *dns.Client, fqdn string) (Records, error) {

	// get authority server for our call
	ns, err := getNameserverAddr(client, d.Domain)
	if err != nil {
		return nil, errors.Wrap(err, "getNameserver")
	}

	var msg dns.Msg

	// set our any query
	msg.SetQuestion(
		dns.Fqdn(fqdn),
		dns.TypeANY,
	)

	reply, _, err := client.Exchange(&msg, ns+":53")
	if err != nil {
		return nil, errors.Wrap(err, "Exchange")
	}

	log.Println(reply.String())

	records := make(Records, 0, len(reply.Answer))

	for idx := range reply.Answer {
		records = append(records, NewRecord(d, reply.Answer[idx], RecordSourceANY))
	}

	for idx := range reply.Extra {
		records = append(records, NewRecord(d, reply.Extra[idx], RecordSourceANY))
	}

	return records, nil
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

// look at existing records and check for any deltas
func (d Domain) CheckDelta(client *dns.Client, records Records) (Records, Records, error) {

	// get authority server for our call
	ns, err := getNameserverAddr(client, d.Domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getNameserver")
	}

	original := make(map[uint32]Record, len(records))
	current := make(map[uint32]Record, 0)

	// loop through records and do a query against the current
	for _, record := range records {

		// set map
		original[record.Hash] = record

		// could perhaps use a buffer for these
		var msg dns.Msg

		msg.SetQuestion(
			record.Name,
			record.RRType,
		)

		reply, _, err := client.Exchange(&msg, ns+":53")
		if err != nil {
			return nil, nil, errors.Wrap(err, "Exchange")
		}

		for idx := range reply.Answer {
			delta := NewRecord(d, reply.Answer[idx], RecordSourceANY)
			current[delta.Hash] = delta
		}
	}

	log.Printf("Current: %d Original: %d", len(current), len(original))

	additions := make(Records, 0)
	for key := range current {
		if _, ok := original[key]; !ok {
			// does not exist means it was added!
			additions = append(additions, current[key])
		}
	}

	removals := make(Records, 0)
	for key := range original {
		if _, ok := current[key]; !ok {
			// does not exist means it was removed!
			removals = append(removals, original[key])
		}
	}

	log.Printf("Additions: %d Removals: %d", len(additions), len(removals))

	return additions, removals, nil
}

// string representation
func (d Domain) String() string {
	return fmt.Sprintf(
		"%s / %s",
		d.Owner.Email,
		d.Domain,
	)
}
