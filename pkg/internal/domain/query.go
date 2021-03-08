package domain

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

var commonRecordTypes = []uint16{
	dns.TypeA,
	dns.TypeAAAA,
	dns.TypeCNAME,
	dns.TypeMX,
	dns.TypeNS,
	dns.TypePTR,
	dns.TypeSRV,
	dns.TypeTXT,
	dns.TypeDNSKEY,
	dns.TypeDS,
	dns.TypeNSEC,
	dns.TypeNSEC3,
	dns.TypeRRSIG,
	dns.TypeAFSDB,
	dns.TypeATMA,
	dns.TypeCAA,
	dns.TypeCERT,
	dns.TypeDHCID,
	dns.TypeDNAME,
	dns.TypeHINFO,
}

func (d Domain) QueryEnumerate(client *dns.Client, targets []string) (Records, error) {
	// get authority server for our call
	nameservers, err := getNameservers(client, d.Domain)
	if err != nil {
		return nil, errors.WithMessage(err, "getNameservers")
	}

	records := make(Records, 0)

	maxWorkers := runtime.GOMAXPROCS(0)
	sem := semaphore.NewWeighted(int64(maxWorkers))
	ctx := context.TODO()

	var mtx sync.Mutex
	var g errgroup.Group

	for idx := range targets {
		for _, typ := range commonRecordTypes {
			if err := sem.Acquire(ctx, 1); err != nil {
				return nil, errors.WithMessage(err, "Acquire")
			}

			func(target string, typ uint16) {
				g.Go(func() error {
					defer sem.Release(1)

					var msg dns.Msg

					fqdn := dns.Fqdn(fmt.Sprintf("%s.%s", targets[idx], d.Domain))
					if len(targets[idx]) == 0 {
						fqdn = dns.Fqdn(d.Domain)
					}

					// set our any query
					msg.SetQuestion(
						fqdn,
						typ,
					)

					log.Printf("\t%s", msg.Question[0].String())

					reply, err := query(client, &msg, nameservers)
					if err != nil {
						return errors.WithMessage(err, "query")
					}

					for idx := range reply.Answer {
						r := NewRecord(d, reply.Answer[idx], RecordSourceEnum)
						if r.Fields == "RFC8482" {
							continue
						}
						mtx.Lock()
						records = append(records, r)
						mtx.Unlock()
						if strings.Contains(r.Fields, d.Domain) {
							if len(strings.Fields(r.Fields)) == 1 {
								log.Printf("Found additional target: %s", r.Fields)
								mtx.Lock()
								targets = append(targets, strings.Replace(r.Fields, "."+d.Domain, "", -1))
								mtx.Unlock()
							}
						}
					}

					for idx := range reply.Extra {
						header := reply.Extra[idx].Header()
						if header.Name == "." && header.Rrtype == dns.TypeOPT {
							// EDNS
							continue
						}

						r := NewRecord(d, reply.Extra[idx], RecordSourceEnum)
						if r.Fields == "RFC8482" {
							continue
						}

						mtx.Lock()
						records = append(records, r)
						mtx.Unlock()
						if strings.Contains(r.Fields, d.Domain) {
							if len(strings.Fields(r.Fields)) == 1 {
								log.Printf("Found additional target: %s", r.Fields)
								mtx.Lock()
								targets = append(targets, strings.Replace(r.Fields, "."+d.Domain, "", -1))
								mtx.Unlock()
							}
						}
					}

					return nil
				})
			}(targets[idx], typ)
		}
	}

	if err := g.Wait(); err != nil {
		return nil, errors.WithMessage(err, "Wait")
	}

	return records, nil
}

// perform an any query
func (d Domain) QueryANY(client *dns.Client, fqdn string) (Records, error) {
	// get authority server for our call
	nameservers, err := getNameservers(client, d.Domain)
	if err != nil {
		return nil, errors.WithMessage(err, "getNameservers")
	}

	var msg dns.Msg

	// set our any query
	msg.SetQuestion(
		dns.Fqdn(fqdn),
		dns.TypeANY,
	)

	reply, err := query(client, &msg, nameservers)
	if err != nil {
		return nil, errors.WithMessage(err, "query")
	}

	records := make(Records, 0, len(reply.Answer))

	for idx := range reply.Answer {
		records = append(records, NewRecord(d, reply.Answer[idx], RecordSourceANY))
	}

	for idx := range reply.Extra {
		header := reply.Extra[idx].Header()
		if header.Name == "." && header.Rrtype == dns.TypeOPT {
			// EDNS
			continue
		}

		records = append(records, NewRecord(d, reply.Extra[idx], RecordSourceANY))
	}

	return records, nil
}

func query(client *dns.Client, original *dns.Msg, nameservers []string) (*dns.Msg, error) {

	// not intrested in recursion?
	original.RecursionDesired = false

	// resets
	client.Net = ""

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
		}

		if reply.Truncated {
			// retry
			continue
		}

		return reply, nil
	}

	return nil, errors.New("no query available")
}
