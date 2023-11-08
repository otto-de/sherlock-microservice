package test

import (
	"context"

	"cloud.google.com/go/errorreporting/apiv1beta1/errorreportingpb"
)

var (
	_ errorreportingpb.ReportErrorsServiceServer = &fakeErrorreportingServer{}
)

type fakeErrorreportingServer struct {
	f ReportErrorFunc
}

func (s *fakeErrorreportingServer) ReportErrorEvent(ctx context.Context, req *errorreportingpb.ReportErrorEventRequest) (*errorreportingpb.ReportErrorEventResponse, error) {
	if s.f != nil {
		return s.f(ctx, req)
	}
	return &errorreportingpb.ReportErrorEventResponse{}, nil
}
