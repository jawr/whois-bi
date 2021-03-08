package domain

import (
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/pkg/errors"
)

var whoisTimestampLayouts = map[string]string{
	"ax": "02.01.2006",
	"im": "02/01/2006 15:04:05",
	"mw": "02.01.2006 15:04:05",
	"at": "20060102 15:04:05",
	"is": "Jan 02 2006",
}

func parseWhoisTimestamp(domain string, t string) (time.Time, error) {
	tstamp, err := dateparse.ParseStrict(t)
	if err == nil {
		return tstamp, nil
	}

	parts := strings.Split(strings.ToLower(domain), ".")

	var lastLayout string

	for i := len(parts) - 1; i >= 0; i-- {
		layout, ok := whoisTimestampLayouts[strings.Join(parts[i:], ".")]
		if ok {
			lastLayout = layout
			continue
		}

		if len(lastLayout) > 0 {
			// nothing bigger found
			break
		}
	}

	if len(lastLayout) > 0 {
		tstamp, err = time.Parse(lastLayout, t)
		if err == nil {
			return tstamp, nil
		}
	}

	return time.Time{}, errors.New("No timestamp format found")
}
