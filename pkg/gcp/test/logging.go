package test

import (
	"context"
	"testing"

	"cloud.google.com/go/logging"
	"github.com/abcxyz/lumberjack/clients/go/pkg/testutil"
	"google.golang.org/api/option"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/grpc"
)

type loggingFakeServer struct {
	logpb.UnimplementedLoggingServiceV2Server

	resp      *logpb.WriteLogEntriesResponse
	returnErr error
}

func (s *loggingFakeServer) WriteLogEntries(_ context.Context, req *logpb.WriteLogEntriesRequest) (*logpb.WriteLogEntriesResponse, error) {
	return s.resp, s.returnErr
}

func newTestLoggingClient(t *testing.T) *logging.Client {
	server := &loggingFakeServer{
		resp: &logpb.WriteLogEntriesResponse{},
	}
	// Setup fake Cloud Logging server.
	addr, conn := testutil.TestFakeGRPCServer(t, func(s *grpc.Server) {
		logpb.RegisterLoggingServiceV2Server(s, server)
	})

	c, err := logging.NewClient(
		context.Background(),
		"test-project",
		option.WithGRPCConn(conn),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
	)
	if err != nil {
		t.Fatalf("Logging Client creation with port `%s` failed: %s", addr, err)
	}
	return c
}
