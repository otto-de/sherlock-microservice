package test

import (
	"context"
	"net"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/errorreporting/apiv1beta1/errorreportingpb"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/apiv2/loggingpb"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/otto-de/sherlock-microservice/pkg/gcp"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	gcp.Services
	PubSub       *pubsub.Client
	PubSubServer *pstest.Server

	listener     net.Listener
	grpcServer   *grpc.Server
	psServerConn *grpc.ClientConn
	fes          *fakeErrorreportingServer
}

type ReportErrorFunc func(ctx context.Context, req *errorreportingpb.ReportErrorEventRequest) (*errorreportingpb.ReportErrorEventResponse, error)

type testServicesOption struct {
	f ReportErrorFunc
}

func WithErrorReportChannel(evs chan<- *errorreportingpb.ReportedErrorEvent) testServicesOption {
	return testServicesOption{
		f: func(ctx context.Context, req *errorreportingpb.ReportErrorEventRequest) (*errorreportingpb.ReportErrorEventResponse, error) {
			evs <- req.Event
			return &errorreportingpb.ReportErrorEventResponse{}, nil
		},
	}
}

func WithErrorReportCallback(f ReportErrorFunc) testServicesOption {
	return testServicesOption{
		f: f,
	}
}

func MustMakeTestServices(ctx context.Context, project, serviceName string, opts ...testServicesOption) *server {

	var fes *fakeErrorreportingServer
	for _, opt := range opts {
		if opt.f != nil {
			fes = &fakeErrorreportingServer{
				f: opt.f,
			}
		}
	}

	if fes == nil {
		fes = &fakeErrorreportingServer{}
	}

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	psServer := pstest.NewServer()

	loggingServer := &loggingFakeServer{
		resp: &loggingpb.WriteLogEntriesResponse{},
	}

	grpcServer := grpc.NewServer()
	errorreportingpb.RegisterReportErrorsServiceServer(grpcServer, fes)
	loggingpb.RegisterLoggingServiceV2Server(grpcServer, loggingServer)
	fakeServerAddr := l.Addr().String()

	psServerConn, err := grpc.Dial(psServer.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	pubSub, err := pubsub.NewClient(ctx, project, option.WithGRPCConn(psServerConn))
	if err != nil {
		panic(err)
	}

	errRep, err := errorreporting.NewClient(
		ctx,
		project,
		errorreporting.Config{
			ServiceName: serviceName,
		},
		option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		panic(err)
	}

	lc, err := logging.NewClient(
		ctx,
		project,
		option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(trace.WithSampler(trace.NeverSample()))

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			panic(err)
		}
	}()

	return &server{
		Services: gcp.Services{
			ErrorReporting: errRep,
			Logging:        lc,
			TracerProvider: tp,
		},
		PubSub:       pubSub,
		PubSubServer: psServer,
		fes:          fes,
		grpcServer:   grpcServer,
		listener:     l,
		psServerConn: psServerConn,
	}
}

func (s *server) Close() error {
	// Ignore close errors because usually
	// we are not that particular about testing
	s.PubSub.Close()
	s.Services.Close()
	s.grpcServer.GracefulStop()
	s.PubSubServer.Close()
	s.psServerConn.Close()
	s.listener.Close()
	return nil
}
