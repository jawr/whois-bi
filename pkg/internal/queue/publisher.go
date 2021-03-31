package queue

import (
	"context"
)

type Publisher interface {
	// JSON encode msg and push it to queue
	Run(ctx context.Context) error
	Publish(ctx context.Context, queue string, msg interface{}) error
}
