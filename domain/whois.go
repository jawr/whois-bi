package domain

import (
	"crypto/sha256"
	"time"

	"github.com/go-pg/pg"
	whoisparser "github.com/jawr/whois-parser-go"
	"github.com/likexian/whois-go"
	"github.com/pkg/errors"
)

type Whois struct {
	ID int `sql:",pk"`

	// parent data
	DomainID int    `sql:",notnull"`
	Domain   Domain `pg:"fk:domain_id"`

	Raw []byte `sql:",notnull"`

	Version []byte `sql:",notnull,unique"`

	CreatedDate    time.Time `sql:",notnull"`
	UpdatedDate    time.Time `sql:",notnull"`
	ExpirationDate time.Time `sql:",notnull"`

	// meta data
	AddedAt   time.Time `sql:",notnull,default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`
}

// do a whois lookup and parse the results
func NewWhois(domain Domain) (Whois, error) {
	raw, err := whois.Whois(domain.Domain)
	if err != nil {
		return Whois{}, errors.Wrap(err, "Whois")
	}

	parsed, err := whoisparser.Parse(raw)
	if err != nil {
		return Whois{}, errors.Wrap(err, "Parse")
	}

	// create our Version
	h := sha256.New()

	h.Write([]byte(domain.Domain))

	if parsed.Domain.UpdatedDate.IsZero() {
		// no updated date lets use the entire raw
		h.Write([]byte(raw))
	} else {
		h.Write([]byte(parsed.Domain.UpdatedDate.String()))
	}

	w := Whois{
		DomainID: domain.ID,
		Domain:   domain,

		Raw: []byte(raw),

		Version: h.Sum(nil),

		CreatedDate:    parsed.Domain.CreatedDate,
		UpdatedDate:    parsed.Domain.UpdatedDate,
		ExpirationDate: parsed.Domain.ExpirationDate,
	}

	return w, nil
}

// insert a whois record
func (w *Whois) Insert(db *pg.DB) error {
	_, err := db.Model(w).Returning("*").OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
