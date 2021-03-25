package domain

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/pkg/errors"
)

type Whois struct {
	ID int `pg:",pk" json:"id"`

	// parent data
	DomainID int    `pg:",notnull" json:"domain_id"`
	Domain   Domain `pg:"fk:domain_id,rel:has-one" json:"-"`

	Raw []byte `pg:",notnull" json:"raw"`

	Version []byte `pg:",notnull,unique" json:"version"`

	CreatedDate    time.Time `pg:",notnull" json:"created_date"`
	UpdatedDate    time.Time `pg:",notnull" json:"updated_date"`
	ExpirationDate time.Time `pg:",notnull" json:"expiration_date"`

	DateErrors []string `pg:",notnull" json:"date_errors"`

	// meta data
	AddedAt   time.Time `pg:",notnull,default:now()" json:"added_at"`
	DeletedAt time.Time `pg:",soft_delete" json:"deleted_at"`
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
