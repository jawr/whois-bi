package rabbit

import (
	"bytes"
	"context"
	"log"
	"sync"
	"time"

	"github.com/furdarius/rabbitroutine"
	"github.com/jawr/whois-bi/pkg/internal/queue"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"
)

type Consumer struct {
	name  string
	queue string
	addr  string

	handler queue.ConsumerHandler

	conn *rabbitroutine.Connector

	buffers sync.Pool
}

func NewConsumer(name, queue, addr string) *Consumer {
	config := rabbitroutine.Config{
		Wait: time.Second * 5,
	}

	conn := rabbitroutine.NewConnector(config)

	buffers := sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	return &Consumer{
		name:    name,
		queue:   queue,
		addr:    addr,
		conn:    conn,
		buffers: buffers,
	}
}

func (c *Consumer) Run(ctx context.Context, handler queue.ConsumerHandler) error {
	var wg errgroup.Group

	c.handler = handler

	wg.Go(func() error {
		return c.conn.Dial(ctx, c.addr)
	})

	wg.Go(func() error {
		return c.conn.StartConsumer(ctx, c)
	})

	return wg.Wait()
}

func (c *Consumer) Declare(ctx context.Context, ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		c.queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) Consume(ctx context.Context, ch *amqp.Channel) error {
	msgs, err := ch.Consume(
		c.queue, // queue
		c.name,  // consumer name
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	log.Printf("Start consuming from %s", c.queue)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-msgs:
			if !ok {
				return amqp.ErrClosed
			}

			c.handler(ctx, msg.Body)

			// do we care about errors here
			msg.Ack(false)
		}
	}
	return nil
}
