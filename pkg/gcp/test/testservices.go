package test

import (
	"testing"

	"github.com/otto-de/sherlock-microservice/pkg/gcp"
)

type TestServices struct {
	gcp.Services
	TestErrorReportingService *FakeErrorsService
}

func NewTestServices(t *testing.T, serviceName string) *TestServices {
	er, es := newTestErrorReportingClient(t, serviceName)
	lc := newTestLoggingClient(t)

	return &TestServices{
		Services: gcp.Services{
			ErrorReporting: er,
			Logging:        lc,
		},
		TestErrorReportingService: es,
	}
}

func (s *TestServices) Close() {
	s.Services.Close()
}
