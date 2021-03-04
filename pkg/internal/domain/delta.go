package domain

import (
	"log"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// look at existing records and check for any deltas
func (d Domain) CheckDelta(client *dns.Client, existing, queried Records) (Records, Records, error) {
	// get authority server for our call
	ns, err := getNameserverAddr(client, d.Domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getNameserver")
	}

	original := make(map[uint32]Record, len(existing))
	current := make(map[uint32]Record, 0)

	// loop through existing records and do a query against the current
	for _, record := range existing {

		// add to original
		original[record.Hash] = record

		var msg dns.Msg

		msg.SetQuestion(
			record.Name,
			record.RRType.V,
		)

		reply, err := query(client, &msg, ns)
		if err != nil {
			return nil, nil, errors.Wrap(err, "query")
		}

		for idx := range reply.Answer {
			delta := NewRecord(d, reply.Answer[idx], RecordSourceANY)
			if delta.Fields == "RFC8482" {
				continue
			}
			current[delta.Hash] = delta
		}
	}

	// loop through queried records and add any to current that do not exist
	for _, record := range queried {
		current[record.Hash] = record
	}

	log.Printf("Current: %d Original: %d", len(current), len(original))

	additions := make(Records, 0)
	for key := range current {
		if _, ok := original[key]; !ok {
			// does not exist means it was added!
			additions = append(additions, current[key])
		}
	}

	removals := make(Records, 0)
	for key := range original {
		if _, ok := current[key]; !ok {
			// does not exist means it was removed!
			removals = append(removals, original[key])
		}
	}

	log.Printf("Additions: %d Removals: %d", len(additions), len(removals))

	return additions, removals, nil
}
