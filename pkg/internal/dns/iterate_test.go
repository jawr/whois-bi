package dns

import (
	"testing"

	"github.com/jawr/whois-bi/pkg/internal/domain"
)

func Test_queryIterate(t *testing.T) {
	c := NewDNSClient()

	type tcase struct {
		name    string
		records []string
	}

	cases := []tcase{
		tcase{
			name: "lawrence.pm",
			records: []string{
				`lawrence.pm.		10799	IN	NS	ns-113-a.gandi.net.`,
				`lawrence.pm.		10799	IN	NS	ns-143-b.gandi.net.`,
				`lawrence.pm.		10799	IN	NS	ns-6-c.gandi.net.`,
				`lawrence.pm.		10799	IN	MX	10 ehlo.mx.ax.`,
				`lawrence.pm.		10799	IN	MX	20 helo.mx.ax.`,
				`lawrence.pm.		10799	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`,
			},
		},
		tcase{
			name: "jl.lu",
			records: []string{
				`jl.lu.			10799	IN	A	116.203.149.40`,
				`jl.lu.			10799	IN	NS	ns-147-a.gandi.net.`,
				`jl.lu.			10799	IN	NS	ns-208-b.gandi.net.`,
				`jl.lu.			10799	IN	NS	ns-112-c.gandi.net.`,
				`*.k3s.jl.lu.		10799	IN	CNAME	traefik.jl.lu.`,
				`traefik.jl.lu.		10799	IN	A	116.203.149.40`,
			},
		},
	}

	targets := []string{
		"",
		// added to target wildcard subdomains and following cnames pointed
		// to the same domain
		"*.k3s",
		"foo.k3s",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			ns, err := c.getNameservers(tc.name)
			if err != nil {
				tt.Fatalf("getNameservers unexpected error: %q", err)
			}

			dom := domain.Domain{Domain: tc.name}

			expectedRecords := domain.Records{}
			for _, er := range tc.records {
				expectedRecords = append(
					expectedRecords,
					domain.NewRecord(
						dom,
						mustCreateRR(tt, er),
						domain.RecordSourceIterate,
					),
				)
			}

			got, err := c.queryIterate(dom, ns, targets)
			if err != nil {
				tt.Fatalf("queryIterate unexpected error: %q", err)
			}

			compareRecords(tt, got, expectedRecords)
		})
	}
}
