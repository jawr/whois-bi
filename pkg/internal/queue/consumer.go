package queue

import (
	"context"
	"encoding/json"
)

type ConsumerHandler func(ctx context.Context, b []byte)

// Consumer to a queue and consume raw messages from it
type Consumer interface {
	// returns a channel where undecoded messages are pushed
	Run(ctx context.Context, handler ConsumerHandler) error
}

type MemoryConsumer struct {
	queue chan []byte
}

func NewMemoryConsumer() *MemoryConsumer {
	return &MemoryConsumer{
		queue: make(chan []byte),
	}
}

func (m *MemoryConsumer) Run(ctx context.Context, handler ConsumerHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b := <-m.queue:
			handler(ctx, b)
		}
	}

	return nil
}

func (m *MemoryConsumer) Publish(body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	m.queue <- b

	return nil
}
