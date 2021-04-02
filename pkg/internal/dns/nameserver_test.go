package dns

import "testing"

func compareSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aa := range a {
		var exists bool
		for _, bb := range b {
			if aa == bb {
				exists = true
				break
			}
		}

		if !exists {
			return false
		}
	}

	return true
}

func Test_getNameservers(t *testing.T) {
	type tcase struct {
		domain   string
		expected []string
	}

	cases := []tcase{
		tcase{
			"mx.ax",
			[]string{
				"angela.ns.cloudflare.com",
				"woz.ns.cloudflare.com",
			},
		},
		tcase{
			"jl.lu",
			[]string{
				"ns-147-a.gandi.net",
				"ns-112-c.gandi.net",
				"ns-208-b.gandi.net",
			},
		},
		tcase{
			"lawrence.pm",
			[]string{
				"ns-6-c.gandi.net",
				"ns-6-c.gandi.net",
				"ns-6-c.gandi.net",
			},
		},
	}

	c := NewDNSClient()

	for _, tc := range cases {
		t.Run(tc.domain, func(t *testing.T) {
			got, err := c.getNameservers(tc.domain)
			if err != nil {
				t.Errorf("expected nil, got %q", err)
			}
			if !compareSlice(tc.expected, got) {
				t.Errorf("expected %q got %q", tc.expected, got)
			}
		})
	}
}
