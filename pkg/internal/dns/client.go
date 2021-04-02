package dns

import (
	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/miekg/dns"
)

type Client interface {
	// GetLive checks to see if the provided stored records still exist as well
	// as checking against our list of domains
	GetLive(dom domain.Domain, stored domain.Records) (domain.Records, error)
}

type DNSClient struct {
	dns.Client
}

func NewDNSClient() *DNSClient {
	client := dns.Client{}

	dc := DNSClient{
		Client: client,
	}

	return &dc
}
