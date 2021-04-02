package dns

import (
	"fmt"
	"log"
	"strings"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

var (
	commonRecordTypes = []uint16{
		dns.TypeA,
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

	for tar := range targets {
		for _, typ := range commonRecordTypes {
			var msg dns.Msg

			fqdn := dns.Fqdn(fmt.Sprintf("%s.%s", tar, dom.Domain))
			if len(tar) == 0 {
				fqdn = dns.Fqdn(dom.Domain)
			}

			// set our any query
			msg.SetQuestion(
				fqdn,
				typ,
			)

			reply, err := c.query(&msg, nameservers)
			if err != nil {
				return nil, errors.WithMessage(err, "query")
			}

			for idx := range reply.Answer {
				r := domain.NewRecord(
					dom,
					reply.Answer[idx],
					domain.RecordSourceEnum,
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
					domain.RecordSourceEnum,
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
		}
	}

	return records, nil
}

// perform an any query
func (c *DNSClient) queryANY(dom domain.Domain, nameservers []string) (domain.Records, error) {
	var msg dns.Msg

	// set our any query
	msg.SetQuestion(
		dns.Fqdn(dom.Domain),
		dns.TypeANY,
	)

	reply, err := c.query(&msg, nameservers)
	if err != nil {
		return nil, errors.WithMessage(err, "query")
	}

	records := make(domain.Records, 0, len(reply.Answer))

	for idx := range reply.Answer {
		records = append(
			records,
			domain.NewRecord(
				dom,
				reply.Answer[idx],
				domain.RecordSourceANY,
			),
		)
	}

	for idx := range reply.Extra {
		header := reply.Extra[idx].Header()
		if header.Name == "." && header.Rrtype == dns.TypeOPT {
			// EDNS
			continue
		}

		records = append(
			records,
			domain.NewRecord(
				dom,
				reply.Extra[idx],
				domain.RecordSourceANY,
			),
		)
	}

	return records, nil
}

func (c *DNSClient) query(original *dns.Msg, nameservers []string) (*dns.Msg, error) {

	// not intrested in recursion?
	original.RecursionDesired = false

	// resets
	c.Net = ""

	var triedUdp, triedEdns, triedTcp bool

	for _, ns := range nameservers {
		msg := original.Copy()

		if triedUdp && !triedEdns {
			o := new(dns.OPT)
			o.Hdr.Name = "."
			o.Hdr.Rrtype = dns.TypeOPT
			o.SetUDPSize(dns.DefaultMsgSize)
			msg.Extra = append(msg.Extra, o)
			triedEdns = true

		} else if triedUdp && triedEdns && !triedTcp {

			c.Net = "tcp"
			triedTcp = true

		} else if triedUdp && triedEdns && triedTcp {
			return nil, errors.New("failed all methods")

		} else {
			triedUdp = true
		}

		reply, _, err := c.Exchange(msg, ns+":53")
		if err != nil {
			log.Printf("error in Exchange with %s: %s", ns, err)
			continue
		}

		if reply.Truncated {
			// retry
			continue
		}

		return reply, nil
	}

	return nil, errors.New("no query available")
}
