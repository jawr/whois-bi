package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/types"
	"github.com/miekg/dns"
)

type JsonDate struct{ pg.NullTime }

func (t JsonDate) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006/01/02"))), nil
}

func (t *JsonDate) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, `"`)
	if len(s) == 0 {
		return nil
	}
	var err error
	t.Time, err = time.Parse("2006/01/02", s)
	if err != nil {
		return err
	}
	return nil
}

type JsonDateTime struct{ pg.NullTime }

func (t JsonDateTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006/01/02 15:04"))), nil
}

func (t *JsonDateTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, `"`)
	if len(s) == 0 {
		return nil
	}
	var err error
	t.Time, err = time.Parse("2006/01/02 15:04", s)
	if err != nil {
		return err
	}
	return nil
}

type JsonRRType struct{ V uint16 }

func (t JsonRRType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t JsonRRType) String() string {
	return dns.TypeToString[t.V]
}

var _ types.ValueAppender = (*JsonRRType)(nil)

func (t JsonRRType) AppendValue(b []byte, flags int) []byte {
	if flags == 1 {
		b = append(b, '\'')
	}
	b = append(b, []byte(fmt.Sprintf("%d", t.V))...)
	if flags == 1 {
		b = append(b, '\'')
	}
	return b
}

var _ types.ValueScanner = (*JsonRRType)(nil)

func (t *JsonRRType) ScanValue(rd types.Reader, n int) error {
	if n <= 0 {
		t.V = uint16(0)
		return nil
	}

	tmp, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	t2, err := strconv.Atoi(string(tmp))
	if err != nil {
		return err
	}

	t.V = uint16(t2)

	return nil
}

func (t *JsonRRType) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, `"`)
	t.V = dns.StringToType[s]
	return nil
}

type DisplayDomain struct {
	tableName struct{} `sql:"domains,select:domains,alias:domain"`

	Domain

	Records int

	Whois int
}
