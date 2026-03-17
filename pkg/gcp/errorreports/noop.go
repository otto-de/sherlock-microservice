package errorreports

import (
	"context"
	"net"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/errorreporting/apiv1beta1/errorreportingpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewNoopClient creates an errorreporting.Client connected to a local in-process
// gRPC server that discards all reported errors. It is intended for use in
// tests where real error reporting to GCP is not desired.
func NewNoopClient(ctx context.Context, projectID string, srv errorreportingpb.ReportErrorsServiceServer) (*errorreporting.Client, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	gsrv := grpc.NewServer()
	errorreportingpb.RegisterReportErrorsServiceServer(gsrv, srv)
	fakeServerAddr := l.Addr().String()
	go func() {
		if err := gsrv.Serve(l); err != nil {
			panic(err)
		}
	}()

	return errorreporting.NewClient(
		ctx,
		projectID,
		errorreporting.Config{
			ServiceName: "HttpServerTest",
		},
		option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
}
