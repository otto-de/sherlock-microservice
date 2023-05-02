package test

import (
	"context"

	"cloud.google.com/go/errorreporting/apiv1beta1/errorreportingpb"
)

var (
	_ errorreportingpb.ReportErrorsServiceServer = &fakeErrorreportingServer{}
)

type fakeErrorreportingServer struct {
	ReportedEvents []*errorreportingpb.ReportedErrorEvent
}

func (s *fakeErrorreportingServer) ReportErrorEvent(ctx context.Context, req *errorreportingpb.ReportErrorEventRequest) (*errorreportingpb.ReportErrorEventResponse, error) {
	s.ReportedEvents = append(s.ReportedEvents, req.Event)
	return &errorreportingpb.ReportErrorEventResponse{}, nil
}
