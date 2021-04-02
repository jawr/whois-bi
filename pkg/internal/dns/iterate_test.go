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
				`lawrence.pm.		10799	IN	SOA	ns1.gandi.net. hostmaster.gandi.net. 1617235200 10800 3600 604800 10800`,
				`lawrence.pm.		10799	IN	MX	10 ehlo.mx.ax.`,
				`lawrence.pm.		10799	IN	MX	20 helo.mx.ax.`,
				`lawrence.pm.		10799	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`,
			},
		},
	}

	targets := map[string]struct{}{
		"": struct{}{},
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
