package worker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jawr/whois-bi/pkg/internal/domain"
	"github.com/jawr/whois-bi/pkg/internal/job"
	"github.com/jawr/whois-bi/pkg/internal/queue"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type mockDnsClient struct {
	live domain.Records
	err  error
}

func (c *mockDnsClient) GetLive(dom domain.Domain, stored domain.Records) (domain.Records, error) {
	return c.live, c.err
}

// MustCreateRR returns a dns.RR, failing the test if any errors are encountered
func mustCreateRR(t *testing.T, raw string) dns.RR {
	t.Helper()

	rr, err := dns.NewRR(raw)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	return rr
}

// Create a new Test Worker that uses a mock DNS Client and in memory queues
func createNewWorker() *Worker {
	dnsClient := &mockDnsClient{}
	publisher := queue.NewMemoryPublisher()
	consumer := queue.NewMemoryConsumer()
	return NewWorker(dnsClient, publisher, consumer)
}

func createDomain() domain.Domain {
	return domain.Domain{
		ID:      1,
		Domain:  "whois.bi",
		OwnerID: 1,
	}
}

func createJob() job.Job {
	dom := createDomain()
	return job.Job{
		ID:        1,
		DomainID:  dom.ID,
		Domain:    dom,
		Errors:    []string{},
		CreatedAt: time.Now(),
	}
}

func Test_RunPublishConsumeStop(t *testing.T) {
	w := createNewWorker()

	ctx, cancel := context.WithCancel(context.Background())

	var wg errgroup.Group

	wg.Go(func() error {
		return w.Run(ctx)
	})

	// publish a job to the consumer
	j := createJob()

	if err := w.consumer.(*queue.MemoryConsumer).Publish(&j); err != nil {
		t.Fatalf("Publish() expected nil got %s", err)
	}

	// check the response on the publisher
	responseBody := <-w.publisher.(*queue.MemoryPublisher).Channel

	var response job.Job
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatalf("Unmarshal() unexpected error: %s", err)
	}

	// check job is as expected
	if j.ID != response.ID {
		t.Fatalf("expected job and response ID to match, got: %d", response.ID)
	}
	if !j.CreatedAt.Equal(response.CreatedAt) {
		t.Fatalf("expected job and response CreatedAt to match, got: %s vs %s", j.CreatedAt, response.CreatedAt)
	}
	if j.DomainID != response.DomainID {
		t.Fatalf("expected job and response DomainID to match, got: %d", response.DomainID)
	}

	// shutdown and check error
	cancel()

	if err := wg.Wait(); err != context.Canceled {
		t.Fatalf("Wait() expected Canceled, got: %s", err)
	}
}

func Test_RunAdditions(t *testing.T) {
	w := createNewWorker()

	ctx, cancel := context.WithCancel(context.Background())

	var wg errgroup.Group

	wg.Go(func() error {
		return w.Run(ctx)
	})

	j := createJob()

	// setup the mockDnsClient to contain the additions we want
	w.dnsClient.(*mockDnsClient).live = domain.Records{
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "www.whois.bi.	300	IN	CNAME	traefik.jl.lu."), domain.RecordSourceIterate),
	}

	if err := w.consumer.(*queue.MemoryConsumer).Publish(&j); err != nil {
		t.Fatalf("Publish() expected nil got %s", err)
	}

	// check the response on the publisher
	responseBody := <-w.publisher.(*queue.MemoryPublisher).Channel

	var response job.Job
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatalf("Unmarshal() unexpected error: %s", err)
	}

	// check additions is correct
	if len(response.RecordAdditions) != 4 {
		t.Fatalf("Expected RecordAdditions to be 4, got %d", len(response.RecordAdditions))
	}
	if len(response.RecordRemovals) != 0 {
		t.Fatalf("Expected RecordRemoals to be 0, got %d", len(response.RecordRemovals))
	}

	// shutdown and check error
	cancel()

	if err := wg.Wait(); err != context.Canceled {
		t.Fatalf("Wait() expected Canceled, got: %s", err)
	}
}

func Test_RunRemovals(t *testing.T) {
	w := createNewWorker()

	ctx, cancel := context.WithCancel(context.Background())

	var wg errgroup.Group

	wg.Go(func() error {
		return w.Run(ctx)
	})

	j := createJob()

	// setup the mockDnsClient to contain the additions we want
	w.dnsClient.(*mockDnsClient).live = domain.Records{}

	j.CurrentRecords = domain.Records{
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "www.whois.bi.	300	IN	CNAME	traefik.jl.lu."), domain.RecordSourceIterate),
	}

	if err := w.consumer.(*queue.MemoryConsumer).Publish(&j); err != nil {
		t.Fatalf("Publish() expected nil got %s", err)
	}

	// check the response on the publisher
	responseBody := <-w.publisher.(*queue.MemoryPublisher).Channel

	var response job.Job
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatalf("Unmarshal() unexpected error: %s", err)
	}

	// check additions is correct
	if len(response.RecordAdditions) != 0 {
		t.Fatalf("Expected RecordAdditions to be 0, got %d", len(response.RecordAdditions))
	}
	if len(response.RecordRemovals) != 4 {
		t.Fatalf("Expected RecordRemoals to be 4, got %d", len(response.RecordRemovals))
	}

	// shutdown and check error
	cancel()

	if err := wg.Wait(); err != context.Canceled {
		t.Fatalf("Wait() expected Canceled, got: %s", err)
	}
}

func Test_RunMix(t *testing.T) {
	w := createNewWorker()

	ctx, cancel := context.WithCancel(context.Background())

	var wg errgroup.Group

	wg.Go(func() error {
		return w.Run(ctx)
	})

	j := createJob()

	// setup the mockDnsClient to contain the additions we want
	w.dnsClient.(*mockDnsClient).live = domain.Records{
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, `whois.bi.	86400	IN	NS	ns2.he.net.`), domain.RecordSourceIterate),
	}

	j.CurrentRecords = domain.Records{
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	10 ehlo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "whois.bi.	43200	IN	MX	20 helo.mx.ax."), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, `whois.bi.	43200	IN	TXT	"v=spf1 include:spf.mx.ax ~all"`), domain.RecordSourceIterate),
		domain.NewRecord(j.Domain, mustCreateRR(t, "www.whois.bi.	300	IN	CNAME	traefik.jl.lu."), domain.RecordSourceIterate),
	}

	if err := w.consumer.(*queue.MemoryConsumer).Publish(&j); err != nil {
		t.Fatalf("Publish() expected nil got %s", err)
	}

	// check the response on the publisher
	responseBody := <-w.publisher.(*queue.MemoryPublisher).Channel

	var response job.Job
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatalf("Unmarshal() unexpected error: %s", err)
	}

	// check additions is correct
	if len(response.RecordAdditions) != 1 {
		t.Fatalf("Expected RecordAdditions to be 1, got %d", len(response.RecordAdditions))
	}
	if len(response.RecordRemovals) != 1 {
		t.Fatalf("Expected RecordRemoals to be 1, got %d", len(response.RecordRemovals))
	}

	// shutdown and check error
	cancel()

	if err := wg.Wait(); err != context.Canceled {
		t.Fatalf("Wait() expected Canceled, got: %s", err)
	}
}

func Test_RunLiveError(t *testing.T) {
	w := createNewWorker()

	ctx, cancel := context.WithCancel(context.Background())

	var wg errgroup.Group

	wg.Go(func() error {
		return w.Run(ctx)
	})

	j := createJob()

	// setup the mockDnsClient to contain the additions we want
	w.dnsClient.(*mockDnsClient).err = errors.New("unable to query")

	if err := w.consumer.(*queue.MemoryConsumer).Publish(&j); err != nil {
		t.Fatalf("Publish() expected nil got %s", err)
	}

	// check the response on the publisher
	responseBody := <-w.publisher.(*queue.MemoryPublisher).Channel

	var response job.Job
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatalf("Unmarshal() unexpected error: %s", err)
	}

	// check additions is correct
	if len(response.RecordAdditions) != 0 {
		t.Fatalf("Expected RecordAdditions to be 0, got %d", len(response.RecordAdditions))
	}
	if len(response.RecordRemovals) != 0 {
		t.Fatalf("Expected RecordRemoals to be 0, got %d", len(response.RecordRemovals))
	}
	if len(response.Errors) != 1 {
		t.Fatalf("Expected Errors to be len 1, got %d", len(response.Errors))
	}

	if response.Errors[0] != "GetLive: unable to query" {
		t.Fatalf("Expected error to be 'GetLive: unable to query' got %q", response.Errors[0])
	}

	// shutdown and check error
	cancel()

	if err := wg.Wait(); err != context.Canceled {
		t.Fatalf("Wait() expected Canceled, got: %s", err)
	}
}
