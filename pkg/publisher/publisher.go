package publisher

import (
	"context"

	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
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

// ApplyCloudEventsPubSubOrderingKey extracts a OrderingKey from Options and
// stores OrderingKey in Context for CloudEvents PubSub.
func ApplyCloudEventsPubSubOrderingKey(ctx context.Context, opts ...Option) context.Context {
	for _, opt := range opts {
		if opt.orderingKey != "" {
			ctx = cepubsub.WithOrderingKey(ctx, opt.orderingKey)
		}
	}
	return nil
}
