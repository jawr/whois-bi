package list

import (
	"regexp"
	"sync"
	"time"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/user"
	"github.com/pkg/errors"
)

type ListType = string

const (
	Whitelist ListType = "whitelist"
	Blacklist ListType = "blacklist"
)

type List struct {
	ID int `pg:",pk" json:"id"`

	OwnerID int       `pg:",notnull,unique:list_type_domain_rrtype_record_owner_id" json:"owner_id"`
	Owner   user.User `pg:"fk:owner_id" json:"-"`

	ListType ListType `pg:",notnull,type:text,unique:list_type_domain_rrtype_record_owner_id" json:"list_type"`

	// fields to match
	Domain string `json:"domain" pg:",notnull,unique:list_type_domain_rrtype_record_owner_id"`
	RRType string `json:"rr_type" pg:",notnull,unique:list_type_domain_rrtype_record_owner_id"`
	Record string `json:"record" pg:",notnull,unique:list_type_domain_rrtype_record_owner_id"`

	domainMatch *regexp.Regexp
	recordMatch *regexp.Regexp
	rrtypeMatch *regexp.Regexp
	once        sync.Once

	AddedAt   time.Time `pg:",type:timestamptz,notnull,default:now()" json:"added_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func anchor(s string) string {
	return "^" + s + "$"
}

func (l List) Validate() error {
	if len(l.Domain) == 0 || len(l.Record) == 0 || len(l.RRType) == 0 {
		return errors.New("missing fields")
	}

	if l.Domain != "*" {
		if _, err := regexp.Compile(anchor(l.Domain)); err != nil {
			return errors.WithMessage(err, "Domain")
		}
	}

	if l.Record != "*" {
		if _, err := regexp.Compile(anchor(l.Record)); err != nil {
			return errors.WithMessage(err, "Record")
		}
	}

	if l.RRType != "*" {
		if _, err := regexp.Compile(anchor(l.RRType)); err != nil {
			return errors.WithMessage(err, "RRType")
		}
	}

	return nil
}

func (l *List) Match(record *domain.Record) bool {
	if l.RRType == "*" && l.Domain == "*" && l.Record == "*" {
		return true
	}

	var matches int

	// init regexps if not already done
	l.once.Do(func() {
		if l.Domain != "*" {
			l.domainMatch = regexp.MustCompile(anchor(l.Domain))
		}
		if l.Record != "*" {
			l.recordMatch = regexp.MustCompile(anchor(l.Record))
		}
		if l.RRType != "*" {
			l.rrtypeMatch = regexp.MustCompile(anchor(l.RRType))
		}
	})

	if l.Domain == "*" || (l.domainMatch != nil && l.domainMatch.MatchString(record.Domain.Domain)) {
		matches++
	}

	if l.Record == "*" || (l.recordMatch != nil && l.recordMatch.MatchString(record.Name)) {
		matches++
	}

	if l.RRType == "*" || (l.rrtypeMatch != nil && l.rrtypeMatch.MatchString(record.RRType.String())) {
		matches++
	}

	return matches == 3
}
