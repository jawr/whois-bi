package domain

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/miekg/dns"
)

type RecordSource uint16

const (
	RecordSourceANY = iota
	RecordSourceAXFR
	RecordSourceManual
)

type Record struct {
	ID int `sql:",pk"`

	// parent data
	DomainID int    `sql:",notnull"`
	Domain   Domain `pg:"fk:domain_id" json:"-"`

	// how was this record generated
	RecordSource RecordSource `sql:",notnull"`

	// textual representaion of the record
	Raw string `sql:",notnull"`

	// fields part of the record:
	// 		facebook.com.	59	IN	A	`31.13.76.35`
	//		facebook.com.	59	IN	TXT	`0 issue "digicert.com"`
	Fields string `sql:",notnull"`

	Name   string `sql:",notnull"`
	RRType uint16 `sql:",notnull"`
	Class  uint16 `sql:",notnull"`
	TTL    uint32 `sql:",notnull"`

	// this is a hash of the fields data and the ttl for
	// easy change detection
	Hash uint32 `sql:",notnull,unique"`

	// meta data
	AddedAt   time.Time `sql:",notnull,default:now()"`
	RemovedAt time.Time
	DeletedAt time.Time `pg:",soft_delete"`
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
	_, err := db.Model(r).
		Set("removed_at = now()").
		WherePK().
		Where(`"record"."removed_at" IS NULL`).
		Update()
	if err != nil {
		return err
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
		RRType: header.Rrtype,
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

	// hash our fields so we have an easy way to compare
	h := fnv.New32a()
	h.Write([]byte(record.Raw))
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
	rrtype := dns.TypeToString[r.RRType]
	return fmt.Sprintf(
		"Domain: %s / Record: %s / %s / Fields: %s  TTL: %d [%d]",
		r.Domain.Domain,
		r.Name,
		rrtype,
		r.Fields,
		r.TTL,
		r.Hash,
	)
}

// custom marshaller
func (r *Record) MarshalJSON() ([]byte, error) {
	type Alias Record

	// can be null
	var removedAt, deletedAt string
	if !r.RemovedAt.IsZero() {
		removedAt = r.RemovedAt.Format("2006/01/02 15:04")
	}
	if !r.DeletedAt.IsZero() {
		deletedAt = r.DeletedAt.Format("2006/01/02 15:04")
	}

	return json.Marshal(&struct {
		RRType    string
		AddedAt   string
		RemovedAt string
		DeletedAt string
		*Alias
	}{
		RRType:    dns.TypeToString[r.RRType],
		AddedAt:   r.AddedAt.Format("2006/01/02 15:04"),
		RemovedAt: removedAt,
		DeletedAt: deletedAt,
		Alias:     (*Alias)(r),
	})
}

func (r *Record) UnmarshalJSON(data []byte) error {
	type Alias Record

	aux := &struct {
		RRType    string
		AddedAt   string
		RemovedAt string
		DeletedAt string
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error

	r.AddedAt, err = time.Parse("2006/01/02 15:04", aux.AddedAt)
	if err != nil {
		return err
	}
	if len(aux.RemovedAt) > 0 {
		r.RemovedAt, err = time.Parse("2006/01/02 15:04", aux.RemovedAt)
		if err != nil {
			return err
		}
	}
	if len(aux.DeletedAt) > 0 {
		r.DeletedAt, err = time.Parse("2006/01/02 15:04", aux.DeletedAt)
		if err != nil {
			return err
		}
	}

	r.RRType = dns.StringToType[aux.RRType]

	return nil
}
