package dns

import (
	"strings"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

var (
	subdomainsToCheck = map[string]struct{}{
		"":                   struct{}{},
		"www":                struct{}{},
		"mx":                 struct{}{},
		"media":              struct{}{},
		"assets":             struct{}{},
		"dashboard":          struct{}{},
		"api":                struct{}{},
		"cdn":                struct{}{},
		"download":           struct{}{},
		"downloads":          struct{}{},
		"mail":               struct{}{},
		"applytics":          struct{}{},
		"email":              struct{}{},
		"app":                struct{}{},
		"img":                struct{}{},
		"default._domainkey": struct{}{},
	}
)

// look at stored records and check for any deltas
func (c DNSClient) GetLive(dom domain.Domain, stored domain.Records) (domain.Records, error) {
	// get authority server for our queries
	nameservers, err := c.getNameservers(dom.Domain)
	if err != nil {
		return nil, errors.WithMessage(err, "getNameserver")
	}

	// create a list of targets we want to check against
	for _, r := range stored {
		name := strings.Replace(r.Name, dns.Fqdn(dom.Domain), "", -1)
		if _, ok := subdomainsToCheck[name]; !ok {
			subdomainsToCheck[name] = struct{}{}
		}
	}

	// query against targets
	live, err := c.queryANY(dom, nameservers)
	if err != nil {
		return nil, errors.WithMessage(err, "queryANY")
	}

	iterated, err := c.queryIterate(dom, nameservers, subdomainsToCheck)
	if err != nil {
		return nil, errors.WithMessage(err, "queryIterate")
	}

	live = append(live, iterated...)

	return live, nil
}
