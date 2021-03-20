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

func TestWildcardMatch(t *testing.T) {
	l := createList("*", "*", "*")
	rec := createRecord("whois.bi", "www", dns.TypeA)

	if !l.Match(&rec) {
		t.Error("expected a match")
	}
}

func TestRRTypeMatch(t *testing.T) {
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

func TestDomainMatch(t *testing.T) {
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

func TestRecordMatch(t *testing.T) {
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
