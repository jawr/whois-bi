package dns

import (
	"strings"
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
		tcase{
			name: "ovh.com",
			records: []string{
				`ovh.com.		21599	IN	NS	dns.ovh.net.`,
				`ovh.com.		21599	IN	NS	ns.ovh.net.`,
				`ovh.com.		21599	IN	NS	dns200.anycast.me.`,
				`ovh.com.		21599	IN	NS	ns200.anycast.me.`,
				`ovh.com.		21599	IN	NS	ns10.ovh.net.`,
				`ovh.com.		21599	IN	NS	dns10.ovh.net.`,
				`ovh.com.		21599	IN	MX	5 mx2.ovh.net.`,
				`ovh.com.		21599	IN	MX	1 mx1.ovh.net.`,
				`ovh.com.		3599	IN	A	198.27.92.1`,
				`ovh.com.		21599	IN	TXT	"v=spf1 include:spf.mailjet.com include:mx.ovh.com ~all"`,
				`ovh.com.		21599	IN	TXT	"google-site-verification=J3fSHAVfI5uZPzM4rlKtSiBnE5iC0lxi0k2-pn0aM1U"`,
				`ovh.com.		21599	IN	TXT	"google-site-verification=Il9nne-nVT0DAIF9l7jwlycs1fMuu_pWggen5IYZVlA"`,
			},
		},
	}

	targets := map[string]struct{}{
		"": struct{}{},
		// added to target wildcard subdomains and following cnames pointed
		// to the same domain
		"k3s": struct{}{},
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

			for i := 0; i < 10; i++ {
				got, err := c.queryIterate(dom, ns, targets)
				if err != nil {
					if strings.Contains(err.Error(), "i/o timeout") {
						continue
					}
					tt.Fatalf("queryIterate unexpected error: %q", err)
				}

				compareRecords(tt, got, expectedRecords)
			}
		})
	}
}
