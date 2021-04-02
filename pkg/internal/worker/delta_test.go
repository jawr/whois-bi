package worker

import (
	"testing"

	"github.com/jawr/whois-bi/pkg/internal/domain"
)

func Test_deltaAdditions(t *testing.T) {
	dom := createDomain()

	stored := domain.Records{}
	live := domain.Records{
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "www.whois.bi.	300	IN	CNAME	traefik.jl.lu."), domain.RecordSourceIterate),
	}

	additions, removals := delta(stored, live)
	if len(additions) != 4 {
		t.Error("expected additions to be 4")
	}

	if len(removals) != 0 {
		t.Error("expected removals to be 0")
	}
}

func Test_deltaRemovals(t *testing.T) {
	dom := createDomain()

	live := domain.Records{}
	stored := domain.Records{
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "www.whois.bi.	300	IN	CNAME	traefik.jl.lu."), domain.RecordSourceIterate),
	}

	additions, removals := delta(stored, live)
	if len(removals) != 4 {
		t.Error("expected removals to be 4")
	}

	if len(additions) != 0 {
		t.Error("expected additions to be 0")
	}
}

func Test_delta(t *testing.T) {
	dom := createDomain()

	live := domain.Records{
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	30 eelo.mx.ax."), domain.RecordSourceIterate),
	}
	stored := domain.Records{
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(dom, mustCreateRR(t, "www.whois.bi.	300	IN	CNAME	traefik.jl.lu."), domain.RecordSourceIterate),
	}

	additions, removals := delta(stored, live)
	if len(removals) != 2 {
		t.Error("expected removals to be 2")
	}

	if len(additions) != 1 {
		t.Error("expected additions to be 1")
	}
}
