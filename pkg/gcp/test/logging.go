package test

import (
	"context"

	"cloud.google.com/go/logging/apiv2/loggingpb"
)

var (
	_ loggingpb.LoggingServiceV2Server = &loggingFakeServer{}
)

type loggingFakeServer struct {
	loggingpb.UnimplementedLoggingServiceV2Server

	resp      *loggingpb.WriteLogEntriesResponse
	returnErr error
}

func (s *loggingFakeServer) WriteLogEntries(_ context.Context, req *loggingpb.WriteLogEntriesRequest) (*loggingpb.WriteLogEntriesResponse, error) {
	return s.resp, s.returnErr
}
