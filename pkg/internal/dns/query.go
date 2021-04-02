package dns

import (
	"log"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// query against the auth nameservers using UDP, EDNS and TCP
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
