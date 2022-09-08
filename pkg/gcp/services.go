package gcp

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Services contains all Google Cloud Services that
// we use
type Services struct {
	Logging        *logging.Client
	ErrorReporting *errorreporting.Client
	TracerProvider *sdktrace.TracerProvider
}

// NewLoggingClient creates a Client which also handles Errors
// by simply writing to Stderr
func NewLoggingClient(project string) (*logging.Client, error) {
	loggingCtx := context.Background()
	loggingClient, err := logging.NewClient(loggingCtx, project)
	if err != nil {
		return nil, err
	}
	loggingClient.OnError = func(err error) {
		fmt.Fprintf(os.Stderr, "Could not log due to error: %v", err)
	}

	return loggingClient, nil
}

// DiscoverServices builds clients for all Services that we use.
func DiscoverServices(project, serviceName string, tracerProviderOptions []sdktrace.TracerProviderOption) (*Services, error) {

	loggingClient, err := NewLoggingClient(project)
	if err != nil {
		return nil, err
	}
	logger := loggingClient.Logger("ErrorReporting")

	errorReportingCtx := context.Background()
	errorClient, err := errorreporting.NewClient(errorReportingCtx, project, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			logger.Log(logging.Entry{
				Severity: logging.Alert,
				Payload:  fmt.Sprintf("Error reporting failed: %s", err),
			})
			// Ignore err since this probably is due to problems with
			// Permissions or Quotas and thus this needs to be fixed first
		},
	})
	if err != nil {
		panic(err)
	}

	exporter, err := texporter.New(texporter.WithProjectID(project))
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(append(tracerProviderOptions, sdktrace.WithBatcher(exporter))...)

	return &Services{
		Logging:        loggingClient,
		ErrorReporting: errorClient,
		TracerProvider: tp,
	}, nil
}

// Close closes all Clients that were created.
// Does **not** handle errors in close since there usually
// is not much that can be done on Close failure anyway.
func (s *Services) Close() {
	s.TracerProvider.ForceFlush(context.Background()) // flushes any pending spans
	s.ErrorReporting.Close()
	s.Logging.Close()
}
