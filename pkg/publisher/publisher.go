package publisher

import (
	"context"

	"github.com/otto-de/sherlock-microservice/pkg/gcp/errorreports"
)

type Option struct {
	orderingKey string
}

func WithOrderingKey(orderingKey string) Option {
	return Option{
		orderingKey: orderingKey,
	}
}

type Publisher[EV any] interface {
	// Publish publishes an Event to Target.
	Publish(ctx context.Context, event *EV, opts ...Option) errorreports.Error
}
