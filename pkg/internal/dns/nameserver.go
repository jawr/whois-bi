package dns

import (
	"strings"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func (c *DNSClient) getNameservers(domain string) ([]string, error) {
	var msg dns.Msg

	msg.SetQuestion(
		dns.Fqdn(domain),
		dns.TypeNS,
	)

	reply, _, err := c.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		return nil, errors.Wrap(err, "Exchange")
	}

	var nameservers []string

	for idx := range reply.Answer {
		ns, ok := reply.Answer[idx].(*dns.NS)
		if !ok {
			return nil, errors.New("casting ns")
		}

		nameservers = append(
			nameservers,
			strings.TrimSuffix(ns.Ns, "."),
		)
	}

	if len(nameservers) == 0 {
		return nil, errors.New("no nameserver found")
	}

	return nameservers, nil
}
