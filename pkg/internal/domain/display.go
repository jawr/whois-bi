package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg/v10/types"
	"github.com/miekg/dns"
)

type DisplayDomain struct {
	tableName struct{} `pg:"domains,select:domains,alias:domain"`

	Domain

	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`

	Records int `json:"records"`

	Whois int `json:"whois"`
}

type JsonRRType struct{ V uint16 }

func (t JsonRRType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t JsonRRType) String() string {
	return dns.TypeToString[t.V]
}

var _ types.ValueAppender = (*JsonRRType)(nil)

func (t JsonRRType) AppendValue(b []byte, flags int) ([]byte, error) {
	if flags == 1 {
		b = append(b, '\'')
	}
	b = append(b, []byte(fmt.Sprintf("%d", t.V))...)
	if flags == 1 {
		b = append(b, '\'')
	}
	return b, nil
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
