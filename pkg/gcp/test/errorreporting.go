package test

import (
	"context"

	"cloud.google.com/go/errorreporting/apiv1beta1/errorreportingpb"
)

var (
	_ errorreportingpb.ReportErrorsServiceServer = &fakeErrorreportingServer{}
)

type fakeErrorreportingServer struct {
	c              chan<- *errorreportingpb.ReportedErrorEvent
	ReportedEvents []*errorreportingpb.ReportedErrorEvent
}

func (s *fakeErrorreportingServer) ReportErrorEvent(ctx context.Context, req *errorreportingpb.ReportErrorEventRequest) (*errorreportingpb.ReportErrorEventResponse, error) {
	if s.c != nil {
		s.c <- req.Event
	}
	s.ReportedEvents = append(s.ReportedEvents, req.Event)
	return &errorreportingpb.ReportErrorEventResponse{}, nil
}
