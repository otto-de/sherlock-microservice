package gcp

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
)

// Services contains all Google Cloud Services that
// we use
type Services struct {
	Logging        *logging.Client
	ErrorReporting *errorreporting.Client
}

// DiscoverServices builds clients for all Services that we use.
func DiscoverServices(project, serviceName string) (*Services, error) {
	loggingCtx := context.Background()
	loggingClient, err := logging.NewClient(loggingCtx, project)
	if err != nil {
		return nil, err
	}
	loggingClient.OnError = func(err error) {
		fmt.Fprintf(os.Stderr, "Could not log due to error: %v", err)
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

	return &Services{
		Logging:        loggingClient,
		ErrorReporting: errorClient,
	}, nil
}

// Close closes all Clients that were created.
// Does **not** handle errors in close since there usually
// is not much that can be done on Close failure anyway.
func (s *Services) Close() {
	s.Logging.Close()
	s.ErrorReporting.Close()
}
