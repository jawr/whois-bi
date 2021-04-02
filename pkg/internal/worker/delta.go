package worker

import (
	"github.com/jawr/whois-bi/pkg/internal/domain"
)

func delta(stored, live domain.Records) (domain.Records, domain.Records) {
	original := make(map[uint32]domain.Record, len(stored))
	current := make(map[uint32]domain.Record, 0)

	for _, record := range stored {
		original[record.Hash] = record
	}

	for _, record := range live {
		current[record.Hash] = record
	}

	additions := make(domain.Records, 0)
	for key := range current {
		if _, ok := original[key]; !ok {
			// does not exist means it was added!
			additions = append(additions, current[key])
		}
	}

	removals := make(domain.Records, 0)
	for key := range original {
		if _, ok := current[key]; !ok {
			// does not exist means it was removed!
			removals = append(removals, original[key])
		}
	}

	return additions, removals
}
