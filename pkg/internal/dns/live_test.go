package dns

import (
	"testing"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
)

func Test_GetLive(t *testing.T) {
	type tcase struct {
		domain  string
		stored  []string
		records []string
	}

	cases := []tcase{
		tcase{
			domain: "lawrence.pm",
			stored: []string{
				`mxax._domainkey.lawrence.pm. 10799 IN	TXT	"v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyGQTP6NZt4vsUbaiQjkcFagUkR6Y3YnGo8xrOApVgIPxMMpUjHkG5VSE" "jAyw51TzwW7qUkM1n2Ehrb+9reBZGQrWKSxUvrx66YZIVEdMQjNN6g3FhPPAsrp7RFWW9CXSqCx9YFXFnA+FcGGjT0EQhl66xNhO" "VFYNJF0wAE0a7uP4SzZF0gO3ATgxZBeSvQdOSCrSwDOb7cnsdT0Jp72UaOWWtxXw/VtTQDNw7NsH0LLlLSw+b7tLE/hV8fyhrFpW" "tqLgStELpXTnzvJ+8CoUo/iZuJSFBhKcSnun+d3GHSJj4Xg0Vj0P4P/wyiEgcxttVyY52IN7XFtKo0r8B1CyCQIDAQAB"`,
			},
			records: []string{
				`lawrence.pm.		10799	IN	NS	ns-113-a.gandi.net.`,
				`lawrence.pm.		10799	IN	NS	ns-143-b.gandi.net.`,
				`lawrence.pm.		10799	IN	NS	ns-6-c.gandi.net.`,
				`lawrence.pm.		10799	IN	SOA	ns1.gandi.net. hostmaster.gandi.net. 1617235200 10800 3600 604800 10800`,
				`lawrence.pm.		10799	IN	MX	10 ehlo.mx.ax.`,
				`lawrence.pm.		10799	IN	MX	20 helo.mx.ax.`,
				`lawrence.pm.		10799	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`,
				`mxax._domainkey.lawrence.pm. 10799 IN	TXT	"v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyGQTP6NZt4vsUbaiQjkcFagUkR6Y3YnGo8xrOApVgIPxMMpUjHkG5VSE" "jAyw51TzwW7qUkM1n2Ehrb+9reBZGQrWKSxUvrx66YZIVEdMQjNN6g3FhPPAsrp7RFWW9CXSqCx9YFXFnA+FcGGjT0EQhl66xNhO" "VFYNJF0wAE0a7uP4SzZF0gO3ATgxZBeSvQdOSCrSwDOb7cnsdT0Jp72UaOWWtxXw/VtTQDNw7NsH0LLlLSw+b7tLE/hV8fyhrFpW" "tqLgStELpXTnzvJ+8CoUo/iZuJSFBhKcSnun+d3GHSJj4Xg0Vj0P4P/wyiEgcxttVyY52IN7XFtKo0r8B1CyCQIDAQAB"`,
				`_dmarc.lawrence.pm.	10799	IN	TXT	"v=DMARC1; p=quarantine"`,
			},
		},
		tcase{
			domain: "jl.lu",
			stored: []string{
				`*.k3s.jl.lu.		10799	IN	CNAME	traefik.jl.lu.`,
			},
			records: []string{
				`*.k3s.jl.lu.		10799	IN	CNAME	traefik.jl.lu.`,
				`jl.lu.			10799	IN	A	116.203.149.40`,
				`jl.lu.			10799	IN	NS	ns-147-a.gandi.net.`,
				`jl.lu.			10799	IN	NS	ns-208-b.gandi.net.`,
				`jl.lu.			10799	IN	NS	ns-112-c.gandi.net.`,
				`jl.lu.			10799	IN	SOA	ns1.gandi.net. hostmaster.gandi.net. 1617235200 10800 3600 604800 10800`,
				`traefik.jl.lu.		10799	IN	A	116.203.149.40`,
			},
		},
	}

	c := NewDNSClient()

	for _, tc := range cases {
		t.Run(tc.domain, func(tt *testing.T) {
			dom := domain.Domain{Domain: tc.domain}

			stored := domain.Records{}
			for _, r := range tc.stored {
				stored = append(stored, domain.NewRecord(dom, mustCreateRR(tt, r), domain.RecordSourceIterate))
			}

			got, err := c.GetLive(dom, stored)
			if err != nil {
				tt.Fatalf("GetLive() unexpected error: %q", err)
			}

			expected := domain.Records{}
			for _, r := range tc.records {
				expected = append(expected, domain.NewRecord(dom, mustCreateRR(tt, r), domain.RecordSourceIterate))
			}

			compareRecords(tt, got, expected)
		})
	}
}

func compareRecords(t *testing.T, got, expected domain.Records) {
	t.Helper()

	for _, g := range got {
		if g.RRType.V == dns.TypeSOA {
			continue
		}
		var exists bool
		for _, e := range expected {
			if g.Hash == e.Hash {
				exists = true
				break
			}
		}
		if !exists {
			t.Errorf("Unexpected got record: %q", g.Raw)
		}
	}

	for _, g := range expected {
		if g.RRType.V == dns.TypeSOA {
			continue
		}
		var exists bool
		for _, e := range got {
			if g.Hash == e.Hash {
				exists = true
				break
			}
		}
		if !exists {
			t.Errorf("Unexpected expected record: %q", g.Raw)
		}
	}
}

// MustCreateRR returns a dns.RR, failing the test if any errors are encountered
func mustCreateRR(t *testing.T, raw string) dns.RR {
	t.Helper()

	rr, err := dns.NewRR(raw)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	return rr
}
