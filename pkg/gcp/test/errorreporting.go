package test

import (
	"context"
	"net"
	"testing"

	"cloud.google.com/go/errorreporting"
	"google.golang.org/api/option"
	clouderrorreportingpb "google.golang.org/genproto/googleapis/devtools/clouderrorreporting/v1beta1"
	"google.golang.org/grpc"
)

type FakeErrorsService struct {
	Messages []string
}

// newTestErrorReportingClient returns a ErrorReportingClient
func newTestErrorReportingClient(t *testing.T, serviceName string) (*errorreporting.Client, *FakeErrorsService) {
	fakeServer := &FakeErrorsService{}
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	gsrv := grpc.NewServer()
	clouderrorreportingpb.RegisterReportErrorsServiceServer(gsrv, fakeServer)
	fakeServerAddr := l.Addr().String()
	go func() {
		if err := gsrv.Serve(l); err != nil {
			panic(err)
		}
	}()
	ec, err := errorreporting.NewClient(context.Background(), "testproject", errorreporting.Config{
		ServiceName:    serviceName,
		ServiceVersion: "1.0",
	}, option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()))
	if err != nil {
		t.Fatal("Creating ErrorReporting Client failed:", err)
	}
	return ec, fakeServer
}

func (s *FakeErrorsService) ReportErrorEvent(ctx context.Context, req *clouderrorreportingpb.ReportErrorEventRequest) (*clouderrorreportingpb.ReportErrorEventResponse, error) {
	s.Messages = append(s.Messages, req.Event.Message)
	return &clouderrorreportingpb.ReportErrorEventResponse{}, nil
}

func (s *FakeErrorsService) Close() error {
	return nil
}
