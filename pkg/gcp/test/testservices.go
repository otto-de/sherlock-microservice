package test

import (
	"testing"

	"github.com/otto-de/sherlock-microservice/pkg/gcp"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TestServices struct {
	gcp.Services
	TestErrorReportingService *FakeErrorsService
}

func NewTestServices(t *testing.T, serviceName string) *TestServices {
	er, es := newTestErrorReportingClient(t, serviceName)
	lc := newTestLoggingClient(t)
	tp := trace.NewTracerProvider(trace.WithSampler(trace.NeverSample()))

	return &TestServices{
		Services: gcp.Services{
			ErrorReporting: er,
			Logging:        lc,
			TracerProvider: tp,
		},
		TestErrorReportingService: es,
	}
}

func (s *TestServices) Close() {
	s.Services.Close()
}
