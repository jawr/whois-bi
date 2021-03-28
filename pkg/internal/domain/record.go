package domain

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/miekg/dns"
)

type RecordSource uint16

const (
	RecordSourceANY = iota
	RecordSourceAXFR
	RecordSourceManual
	RecordSourceEnum
)

type Record struct {
	ID int `pg:",pk" json:"id"`

	// parent data
	DomainID int    `pg:",notnull" json:"domain_id"`
	Domain   Domain `pg:"fk:domain_id,rel:has-one" json:"-"`

	// how was this record generated
	RecordSource RecordSource `pg:",notnull,use_zero" json:"record_source"`

	// textual representaion of the record
	Raw string `pg:",notnull" json:"raw"`

	// fields part of the record:
	// 		facebook.com.	59	IN	A	`31.13.76.35`
	//		facebook.com.	59	IN	TXT	`0 issue "digicert.com"`
	Fields string `pg:",notnull" json:"fields"`

	Name   string     `pg:",notnull" json:"name"`
	RRType JsonRRType `pg:",notnull" json:"rr_type"`
	Class  uint16     `pg:",notnull" json:"rr_class"`
	TTL    uint32     `pg:",notnull,use_zero" json:"ttl"`

	// this is a hash of the fields data and the ttl for
	// easy change detection
	Hash uint32 `pg:",notnull,unique" json:"hash"`

	// meta data
	AddedAt   time.Time `pg:",type:timestamptz,notnull,default:now()" json:"added_at"`
	RemovedAt time.Time `pg:",type:timestamptz" json:"removed_at"`
	DeletedAt time.Time `pg:",soft_delete" json:"deleted_at"`
}

// helper type
type Records []Record

// insert all records
func (r *Records) Insert(db *pg.DB) error {
	if len(*r) == 0 {
		return nil
	}
	_, err := db.Model(r).
		OnConflict("DO NOTHING").
		Returning("*").
		Insert()
	if err != nil {
		return err
	}
	return nil
}

func (r *Records) Remove(db *pg.DB) error {
	if len(*r) == 0 {
		return nil
	}
	for _, record := range *r {
		_, err := db.Model(&record).
			Set("removed_at = now()").
			WherePK().
			Where(`"record"."removed_at" IS NULL`).
			Update()
		if err != nil {
			return err
		}
	}
	return nil
}

// convert a dns.RR to Record
func NewRecord(domain Domain, rr dns.RR, source RecordSource) Record {

	header := rr.Header()

	record := Record{
		DomainID: domain.ID,
		Domain:   domain,

		RecordSource: source,

		Raw: rr.String(),

		Name:   header.Name,
		RRType: JsonRRType{header.Rrtype},
		Class:  header.Class,
		TTL:    header.Ttl,
	}

	// create and set our fields data
	numFields := dns.NumField(rr)

	fields := make([]string, 0, numFields)

	for i := 0; i <= numFields; i++ {
		field := dns.Field(rr, i)
		if len(field) > 0 {
			fields = append(fields, field)
		}
	}

	record.Fields = strings.Join(fields, " ")

	// hash fields + rrtype + name
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf(
		"%s%d%s",
		record.Name,
		record.RRType,
		record.Fields,
	)))
	record.Hash = h.Sum32()

	return record
}

// insert a record in to the database
func (r *Record) Insert(db *pg.DB) error {
	if _, err := db.Model(r).Returning("*").Insert(); err != nil {
		return err
	}
	return nil
}

// set a record as removed
func (r Record) Remove(db *pg.DB) error {
	_, err := db.Model(&r).
		Set("removed_at = now()").
		Where("id = ? AND removed_at IS NULL", r.ID).
		Update()
	if err != nil {
		return err
	}
	return nil
}

// string representation
func (r Record) String() string {
	return fmt.Sprintf(
		"Domain: %s / Record: %s / %s / Fields: %s  TTL: %d [%d]",
		r.Domain.Domain,
		r.Name,
		r.RRType,
		r.Fields,
		r.TTL,
		r.Hash,
	)
}
