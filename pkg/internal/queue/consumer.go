package queue

import (
	"context"

	"github.com/streadway/amqp"
)

type ConsumerHandler func(ctx context.Context, msg *amqp.Delivery)

// Consumer to a queue and consume raw messages from it
type Consumer interface {
	// returns a channel where undecoded messages are pushed
	Run(ctx context.Context, handler ConsumerHandler) error
}
