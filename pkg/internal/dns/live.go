package dns

import (
	"strings"
	"time"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

var (
	subdomainsToCheck = []string{
		"",
		"*",
		"www",
		"mx",
		"media",
		"assets",
		"dashboard",
		"api",
		"cdn",
		"download",
		"downloads",
		"mail",
		"applytics",
		"email",
		"app",
		"img",
		"default._domainkey",
		"_dmarc",
		"spf",
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
		var exists bool
		for _, t := range subdomainsToCheck {
			if t == name {
				exists = true
				break
			}
		}
		if !exists {
			subdomainsToCheck = append(subdomainsToCheck, name)
		}
	}

	var live domain.Records

	for i := 0; i < 10; i++ {
		live, err = c.queryIterate(dom, nameservers, subdomainsToCheck)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				time.Sleep(time.Millisecond * 500)
				continue
			}
		}
		break
	}

	return live, err
}
