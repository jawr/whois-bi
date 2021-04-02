package queue

import (
	"context"
	"encoding/json"
)

type Publisher interface {
	// JSON encode msg and push it to queue
	Run(ctx context.Context) error
	Publish(ctx context.Context, queue string, msg interface{}) error
}

type MemoryPublisher struct {
	Channel chan []byte
}

func NewMemoryPublisher() *MemoryPublisher {
	return &MemoryPublisher{
		Channel: make(chan []byte),
	}
}

func (m *MemoryPublisher) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (m *MemoryPublisher) Publish(ctx context.Context, queue string, msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	m.Channel <- b
	return nil
}
