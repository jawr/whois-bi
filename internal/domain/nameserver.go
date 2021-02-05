package domain

import (
	"strings"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func getNameserverAddr(client *dns.Client, domain string) (string, error) {
	var msg dns.Msg

	msg.SetQuestion(
		dns.Fqdn(domain),
		dns.TypeNS,
	)

	reply, _, err := client.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		return "", errors.Wrap(err, "Exchange")
	}

	var nameserver string

	for idx := range reply.Answer {
		ns, ok := reply.Answer[idx].(*dns.NS)
		if !ok {
			return "", errors.New("casting ns")
		}

		nameserver = ns.Ns

		break
	}

	if len(nameserver) == 0 {
		return "", errors.New("no nameserver found")
	}

	return strings.TrimSuffix(nameserver, "."), nil
}
