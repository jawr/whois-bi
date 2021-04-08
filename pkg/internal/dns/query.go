package dns

import (
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// query against the auth nameservers using UDP, EDNS and TCP
func (c *DNSClient) query(original *dns.Msg, nameservers []string) (*dns.Msg, error) {
	for i := 0; i < 10; i++ {
		msg, err := c.rawquery(original, nameservers)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				time.Sleep(time.Millisecond * 250)
				continue
			}

			// fall through and return msg and err
		}

		return msg, err
	}

	return nil, errors.New("timeout")
}

func (c *DNSClient) rawquery(original *dns.Msg, nameservers []string) (*dns.Msg, error) {

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
