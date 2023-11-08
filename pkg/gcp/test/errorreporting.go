package test

import (
	"context"
	"testing"

	"cloud.google.com/go/errorreporting/apiv1beta1/errorreportingpb"
)

var (
	_ errorreportingpb.ReportErrorsServiceServer = &fakeErrorreportingServer{}
)

type fakeErrorreportingServer struct {
	tForFailOnEvent *testing.T
	ReportedEvents  []*errorreportingpb.ReportedErrorEvent
}

func (s *fakeErrorreportingServer) ReportErrorEvent(ctx context.Context, req *errorreportingpb.ReportErrorEventRequest) (*errorreportingpb.ReportErrorEventResponse, error) {
	if s.tForFailOnEvent != nil {
		s.tForFailOnEvent.Fatal("Unexpected ErrorEvent reported:", req.Event)
	}
	s.ReportedEvents = append(s.ReportedEvents, req.Event)
	return &errorreportingpb.ReportErrorEventResponse{}, nil
}
