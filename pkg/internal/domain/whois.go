package domain

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/pkg/errors"
)

type Whois struct {
	ID int `sql:",pk"`

	// parent data
	DomainID int    `sql:",notnull"`
	Domain   Domain `pg:"fk:domain_id" json:"-"`

	Raw []byte `sql:",notnull"`

	Version []byte `sql:",notnull,unique"`

	CreatedDate    time.Time `sql:",notnull"`
	UpdatedDate    time.Time `sql:",notnull"`
	ExpirationDate time.Time `sql:",notnull"`

	DateErrors []string `sql:",notnull"`

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

	// parse dates

	var dateErrors []string

	createdDate, err := parseWhoisTimestamp(domain.Domain, parsed.Domain.CreatedDate)
	if err != nil {
		dateErrors = append(dateErrors, fmt.Sprintf("createdDate: %s", err))
	}

	updatedDate, err := parseWhoisTimestamp(domain.Domain, parsed.Domain.UpdatedDate)
	if err != nil {
		dateErrors = append(dateErrors, fmt.Sprintf("updatedDate: %s", err))
	}

	expirationDate, err := parseWhoisTimestamp(domain.Domain, parsed.Domain.ExpirationDate)
	if err != nil {
		dateErrors = append(dateErrors, fmt.Sprintf("expirationDate: %s", err))
	}

	// finish writing our hash
	if updatedDate.IsZero() {
		// no updated date lets use the entire raw
		h.Write([]byte(raw))
	} else {
		h.Write([]byte(updatedDate.String()))
	}

	w := Whois{
		DomainID: domain.ID,
		Domain:   domain,

		Raw: []byte(raw),

		Version: h.Sum(nil),

		CreatedDate:    createdDate,
		UpdatedDate:    updatedDate,
		ExpirationDate: expirationDate,

		DateErrors: dateErrors,
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
