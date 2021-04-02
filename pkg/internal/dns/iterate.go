package dns

import (
	"fmt"
	"strings"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

var (
	commonRecordTypes = []uint16{
		dns.TypeA,
		dns.TypeANY,
		dns.TypeAAAA,
		dns.TypeCNAME,
		dns.TypeMX,
		dns.TypeNS,
		dns.TypePTR,
		dns.TypeSRV,
		dns.TypeTXT,
		dns.TypeDNSKEY,
		dns.TypeDS,
		dns.TypeNSEC,
		dns.TypeNSEC3,
		dns.TypeRRSIG,
		dns.TypeAFSDB,
		dns.TypeATMA,
		dns.TypeCAA,
		dns.TypeCERT,
		dns.TypeDHCID,
		dns.TypeDNAME,
		dns.TypeHINFO,
	}
)

func (c *DNSClient) queryIterate(dom domain.Domain, nameservers []string, targets map[string]struct{}) (domain.Records, error) {
	records := make(domain.Records, 0)

	// currently only handles wildcards with depth of 1 correctly
	wildcards := make(map[uint16]int, 0)

	for tar := range targets {
		tar := strings.TrimSuffix(tar, ".")

		depth := len(strings.Split(tar, "."))

		for _, typ := range commonRecordTypes {
			if tar == "*" {
				if wdepth, ok := wildcards[typ]; ok && wdepth == depth {
					continue
				}
			}

			var msg dns.Msg

			fqdn := dns.Fqdn(fmt.Sprintf("%s.%s", tar, dom.Domain))
			if len(tar) == 0 {
				fqdn = dns.Fqdn(dom.Domain)
			}

			// set our any query
			msg.SetQuestion(fqdn, typ)

			reply, err := c.query(&msg, nameservers)
			if err != nil {
				return nil, errors.WithMessagef(err, "query %q", msg.String())
			}

			if tar == "*" && len(reply.Answer) > 0 {
				wildcards[typ] = depth
			}

			for idx := range reply.Answer {
				r := domain.NewRecord(
					dom,
					reply.Answer[idx],
					domain.RecordSourceIterate,
				)
				if r.Fields == "RFC8482" {
					continue
				}

				records = append(records, r)

				if strings.Contains(r.Fields, dom.Domain) {
					if len(strings.Fields(r.Fields)) == 1 {
						name := strings.Replace(
							r.Fields,
							fmt.Sprintf(".%s", dom.Domain),
							"",
							-1,
						)
						if _, ok := targets[name]; !ok {
							targets[name] = struct{}{}
						}
					}
				}
			}

			for idx := range reply.Extra {
				header := reply.Extra[idx].Header()
				if header.Name == "." && header.Rrtype == dns.TypeOPT {
					// EDNS
					continue
				}

				r := domain.NewRecord(
					dom,
					reply.Extra[idx],
					domain.RecordSourceIterate,
				)

				if r.Fields == "RFC8482" {
					continue
				}

				records = append(records, r)

				// check if we need to append more
				if strings.Contains(r.Fields, dom.Domain) {
					if len(strings.Fields(r.Fields)) == 1 {
						name := strings.Replace(
							r.Fields,
							fmt.Sprintf(".%s", dom.Domain),
							"",
							-1,
						)
						if _, ok := targets[name]; !ok {
							targets[name] = struct{}{}
						}
					}
				}
			}
		}
	}

	return records, nil
}
