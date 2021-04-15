package list

import (
	"testing"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
)

func createList(d, rt, r string) List {
	return List{
		ListType: Blacklist,
		Domain:   d,
		RRType:   rt,
		Record:   r,
	}
}

func createRecord(d, n string, rt uint16) domain.Record {
	return domain.Record{
		Name: n,
		Domain: domain.Domain{
			Domain: d,
		},
		RRType: domain.JsonRRType{rt},
	}
}

func Test_WildcardMatch(t *testing.T) {
	t.Parallel()

	l := createList("*", "*", "*")
	rec := createRecord("whois.bi", "www", dns.TypeA)

	if !l.Match(&rec) {
		t.Error("expected a match")
	}
}

func Test_RRTypeMatch(t *testing.T) {
	t.Parallel()

	l := createList("*", "A", "*")

	passes := []domain.Record{
		createRecord("whois.bi", "www", dns.TypeA),
		createRecord("whois.bi", "mx", dns.TypeA),
		createRecord("mx.ax", "", dns.TypeA),
		createRecord("mx.ax", "@", dns.TypeA),
	}

	fails := []domain.Record{
		createRecord("whois.bi", "www", dns.TypeCNAME),
		createRecord("mx.ax", "@", dns.TypeTXT),
		createRecord("mx.ax", "", dns.TypeSOA),
	}

	for _, r := range fails {
		if l.Match(&r) {
			t.Errorf("expected no match: %s", r.String())
		}
	}

	for _, r := range passes {
		if !l.Match(&r) {
			t.Errorf("expected a match: %s", r.String())
		}
	}
}

func Test_DomainMatch(t *testing.T) {
	t.Parallel()

	l := createList("whois.bi", "*", "*")

	type testcase struct {
		d domain.Domain
		r domain.Record
	}

	passes := []domain.Record{
		createRecord("whois.bi", "www", dns.TypeA),
		createRecord("whois.bi", "mx", dns.TypeCNAME),
		createRecord("whois.bi", "", dns.TypeTXT),
		createRecord("whois.bi", "@", dns.TypeSOA),
	}

	fails := []domain.Record{
		createRecord("mx.ax", "www", dns.TypeCNAME),
		createRecord("mx.ax", "@", dns.TypeTXT),
		createRecord("mx.ax", "", dns.TypeSOA),
	}

	for _, r := range fails {
		if l.Match(&r) {
			t.Errorf("expected no match: %s", r.String())
		}
	}

	for _, r := range passes {
		if !l.Match(&r) {
			t.Errorf("expected a match: %s", r.String())
		}
	}
}

func Test_RecordMatch(t *testing.T) {
	t.Parallel()

	l := createList("*", "*", "www")

	type testcase struct {
		d domain.Domain
		r domain.Record
	}

	passes := []domain.Record{
		createRecord("whois.bi", "www", dns.TypeA),
		createRecord("mx.ax", "www", dns.TypeTXT),
		createRecord("jl.lu", "www", dns.TypeCNAME),
	}

	fails := []domain.Record{
		createRecord("whois.bi", "mx", dns.TypeA),
		createRecord("mx.ax", "mx", dns.TypeTXT),
		createRecord("jl.lu", "mx", dns.TypeCNAME),
	}

	for _, r := range fails {
		if l.Match(&r) {
			t.Errorf("expected no match: %s", r.String())
		}
	}

	for _, r := range passes {
		if !l.Match(&r) {
			t.Errorf("expected a match: %s", r.String())
		}
	}
}

func Test_Validate(t *testing.T) {
	t.Parallel()

	type tcase struct {
		domain   string
		rrtype   string
		record   string
		expected string
	}

	cases := []tcase{
		tcase{
			domain:   "",
			rrtype:   "*",
			record:   "*",
			expected: "missing fields",
		},
		tcase{
			domain:   "*",
			rrtype:   "*",
			record:   "",
			expected: "missing fields",
		},
		tcase{
			domain:   "*",
			rrtype:   "",
			record:   "*",
			expected: "missing fields",
		},
		tcase{
			domain:   "[A--]*",
			rrtype:   "*",
			record:   "*",
			expected: "Domain: error parsing regexp: invalid character class range: `A--`",
		},
		tcase{
			domain:   "*",
			rrtype:   "[A--]*",
			record:   "*",
			expected: "RRType: error parsing regexp: invalid character class range: `A--`",
		},
		tcase{
			domain:   "*",
			rrtype:   "*",
			record:   "[A--]*",
			expected: "Record: error parsing regexp: invalid character class range: `A--`",
		},
	}

	for _, tc := range cases {
		t.Run(tc.expected, func(tt *testing.T) {
			l := List{
				Domain: tc.domain,
				RRType: tc.rrtype,
				Record: tc.record,
			}
			err := l.Validate()
			if err == nil {
				tt.Fatal("Validate() expected an error got nil")
			}
			if err.Error() != tc.expected {
				tt.Fatalf("Validate() expected %q got %q", tc.expected, err)
			}
		})
	}
}

func Test_ValidateSuccess(t *testing.T) {
	t.Parallel()

	l := List{
		Domain: "*",
		RRType: "*",
		Record: "*",
	}

	if err := l.Validate(); err != nil {
		t.Fatalf("Validate() expected nil got %q", err)
	}
}
