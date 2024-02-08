package publisher

import (
	"context"
	"fmt"

	"cloud.google.com/go/errorreporting"
	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/otto-de/sherlock-microservice/pkg/gcp/errorreports"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Option struct {
	extensions  map[string]any
	orderingKey string
}

func WithExtensions(extensions map[string]any) Option {
	return Option{
		extensions: extensions,
	}
}

func WithOrderingKey(orderingKey string) Option {
	return Option{
		orderingKey: orderingKey,
	}
}

type PublishError interface {
	errorreports.Error

	// IsACK returns true if recipient acknowledged the event.
	IsACK() bool
}

type Publisher[EV any] interface {
	// Publish publishes an Event to Target.
	Publish(ctx context.Context, event *EV, opts ...Option) PublishError
}

// ApplyCloudEventsPubSubOrderingKey extracts a OrderingKey from Options and
// stores OrderingKey in Context for CloudEvents PubSub.
func ApplyCloudEventsPubSubOrderingKey(ctx context.Context, opts ...Option) context.Context {
	for _, opt := range opts {
		if opt.orderingKey != "" {
			ctx = cepubsub.WithOrderingKey(ctx, opt.orderingKey)
		}
	}
	return ctx
}

func ApplyCloudEventOptions(ctx context.Context, event *event.Event, opts ...Option) context.Context {
	ctx = ApplyCloudEventsPubSubOrderingKey(ctx, opts...)
	for _, opt := range opts {
		if opt.extensions != nil {
			for k, v := range opt.extensions {
				event.SetExtension(k, v)
			}
		}
	}
	return ctx
}

// PanicPublisher wraps a Publisher.
// Panics if publishing fails.
type PanicPublisher[EV any] struct {
	er   *errorreporting.Client
	base Publisher[EV]
}

func NewPanicPublisher[EV any](er *errorreporting.Client, base Publisher[EV]) *PanicPublisher[EV] {
	return &PanicPublisher[EV]{
		er:   er,
		base: base,
	}
}

func (p *PanicPublisher[EV]) PublishWithNACKPanic(ctx context.Context, event *EV, opts ...Option) errorreports.Error {
	publishErr := p.base.Publish(ctx, event, opts...)
	if publishErr == nil {
		return nil
	}
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, "Send failed")
	span.RecordError(publishErr)
	if publishErr.IsACK() {
		p.er.Report(errorreporting.Entry{
			Error: publishErr,
		})
	} else {
		p.er.ReportSync(ctx, errorreporting.Entry{
			Error: publishErr,
		})
		// For now we do not recover here but just panic
		// TODO: Optimally we should introduce circuit breakers instead
		panicErr := fmt.Errorf("choosing to panic (may trigger restart) due to possible unrecoverable publishing error: %w", publishErr)
		panic(panicErr)
	}
	return publishErr
}
