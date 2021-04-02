package job

import (
	"log"

	"github.com/jawr/whois-bi/pkg/internal/list"
	"github.com/pkg/errors"
)

func (m *Manager) handleLists(response *Job) error {
	// eventually we will want to cache this, but for now its
	// better to give the user a real time feeling as db calls
	// will be cheap
	var whitelists, blacklists []list.List

	if response.Domain.OwnerID == 0 {
		return errors.New("expected a user id")
	}
	uID := response.Domain.OwnerID

	err := m.db.Model(&whitelists).Where("owner_id = ? AND list_type = ?", uID, list.Whitelist).Select()
	if err != nil {
		return err
	}

	err = m.db.Model(&blacklists).Where("owner_id = ? AND list_type = ?", uID, list.Blacklist).Select()
	if err != nil {
		return err
	}

	return handleLists(response, whitelists, blacklists)
}

func handleLists(response *Job, whitelists, blacklists []list.List) error {
	for _, w := range whitelists {
		// if we match anything, drop out before checking blacklists
		for _, r := range response.RecordAdditions {
			if w.Match(&r) {
				return nil
			}
		}
		for _, r := range response.RecordRemovals {
			if w.Match(&r) {
				return nil
			}
		}
	}

	for _, b := range blacklists {
		// if we match anything, delete the record as it
		// is not to be matched

		i := 0
		for _, r := range response.RecordAdditions {
			if !b.Match(&r) {
				response.RecordAdditions[i] = r
				i++
			} else {
				log.Printf("Removing recordAddition %s as matched %d", r.Fields, b.ID)
			}
		}
		response.RecordAdditions = response.RecordAdditions[:i]

		i = 0
		for _, r := range response.RecordRemovals {
			if !b.Match(&r) {
				response.RecordRemovals[i] = r
				i++
			} else {
				log.Printf("Removing recordRemoval %s as matched %d", r.Fields, b.ID)
			}
		}
		response.RecordRemovals = response.RecordRemovals[:i]
	}

	return nil
}
