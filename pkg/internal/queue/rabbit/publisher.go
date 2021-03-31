package rabbit

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/furdarius/rabbitroutine"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"
)

type Publisher struct {
	addr string

	conn      *rabbitroutine.Connector
	publisher *rabbitroutine.FireForgetPublisher

	buffers sync.Pool
}

// Create a new MQ publisher
func NewPublisher(addr string) *Publisher {
	config := rabbitroutine.Config{
		Wait: time.Second * 5,
	}

	conn := rabbitroutine.NewConnector(config)
	pool := rabbitroutine.NewLightningPool(conn)
	publisher := rabbitroutine.NewFireForgetPublisher(pool)

	buffers := sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	p := Publisher{
		addr:      addr,
		conn:      conn,
		publisher: publisher,
		buffers:   buffers,
	}

	return &p
}

// Run long running functions
func (p *Publisher) Run(ctx context.Context) error {
	var wg errgroup.Group

	// run dialer
	wg.Go(func() error {
		return p.conn.Dial(ctx, p.addr)
	})

	return wg.Wait()
}

// Encode and publish a message to provided queue within context
func (p *Publisher) Publish(ctx context.Context, queue string, msg interface{}) error {
	buff := p.buffers.Get().(*bytes.Buffer)
	buff.Reset()
	defer p.buffers.Put(buff)

	if err := json.NewEncoder(buff).Encode(msg); err != nil {
		return err
	}

	err := p.publisher.Publish(ctx, "", queue, amqp.Publishing{
		Body:        buff.Bytes(),
		ContentType: "application/json",
		Timestamp:   time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
