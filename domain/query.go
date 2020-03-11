package domain

import (
	"log"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

func query(client *dns.Client, original *dns.Msg, ns string) (*dns.Msg, error) {
	// resets
	client.Net = ""

	var triedUdp, triedEdns, triedTcp bool

	for {
		msg := original.Copy()

		if triedUdp && !triedEdns {
			log.Println("trying edns")
			o := new(dns.OPT)
			o.Hdr.Name = "."
			o.Hdr.Rrtype = dns.TypeOPT
			o.SetUDPSize(dns.DefaultMsgSize)
			msg.Extra = append(msg.Extra, o)
			triedEdns = true

		} else if triedUdp && triedEdns && !triedTcp {
			log.Println("trying tcp")

			client.Net = "tcp"
			triedTcp = true

		} else if triedUdp && triedEdns && triedTcp {
			return nil, errors.New("failed all methods")

		} else {
			triedUdp = true
		}

		reply, _, err := client.Exchange(msg, ns+":53")
		if err != nil {
			log.Printf("error in Exchange with %s: %s", ns, err)
			continue
			return nil, errors.Wrap(err, "Exchange")
		}

		log.Printf("-------------- /  start\ndns reply \n\n%s\n\n-------------- / end", reply.String())

		if reply.Truncated {
			log.Println("truncated, retrying")
			// retry
			continue
		}

		return reply, nil
	}
}
