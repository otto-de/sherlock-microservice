package gcp

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub/v2"
	"github.com/otto-de/sherlock-microservice/pkg/gke"
	"go.opentelemetry.io/otel/attribute"
	texporter "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

// Services contains all Google Cloud Services that
// we use
// TODO: Spit non-Services to another struct
type Services struct {
	Logging           *logging.Client
	ErrorReporting    *errorreporting.Client
	PubSub            *pubsub.Client
	TracerProvider    *sdktrace.TracerProvider
	MonitoredResource *monitoredres.MonitoredResource
}

// NewLoggingClient creates a Client which also handles Errors
// by simply writing to Stderr
func NewLoggingClient(project string) (*logging.Client, error) {
	loggingCtx := context.Background()
	loggingClient, err := logging.NewClient(loggingCtx, project)
	if err != nil {
		return nil, err
	}
	loggingClient.OnError = func(err error) {
		fmt.Fprintf(os.Stderr, "Could not log due to error: %v", err)
	}

	return loggingClient, nil
}

type discoveryOption struct {
	clusterName             string
	containerName           string
	namespace               string
	pod                     string
	gkeAutoDiscoverMetaData bool
	discoverPubSub          bool
	traceExporterOption     texporter.Option
	traceProviderResource   *resource.Resource
}

func WithKubernetes(clusterName, namespace, pod, containerName string) discoveryOption {
	return discoveryOption{
		clusterName:   clusterName,
		namespace:     namespace,
		pod:           pod,
		containerName: containerName,
	}
}

func WithPubSub() discoveryOption {
	return discoveryOption{
		discoverPubSub: true,
	}
}

func WithGKEAutoDiscoverMetaData() discoveryOption {
	/*
		This option will try to auto discover available metadata from the Google Cloud metadata service and environment variables.
		The metadata will be used to create the monitored resource and trace resource.
		For the container name it will use the environment variable CONTAINER_NAME.
		For the pod name it will use the environment variable POD_NAME.
		For the namespace it will use the environment variable NAMESPACE.
		These variables cant be fetched from the metadata service.
		Set the environment variables in the deployment manifest.

		- name: POD_NAME
			valueFrom:
				fieldRef:
					fieldPath: metadata.name
		- name: POD_NAMESPACE
			valueFrom:
				fieldRef:
					fieldPath: metadata.namespace
		- name: CONTAINER_NAME
			value: test-container
	*/
	return discoveryOption{
		gkeAutoDiscoverMetaData: true,
	}
}

func WithTraceExporterOption(to texporter.Option) discoveryOption {
	return discoveryOption{
		traceExporterOption: to,
	}
}

// Use this instead of https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#WithResource which gets overwritten
func WithTraceProviderResource(tr *resource.Resource) discoveryOption {
	return discoveryOption{
		traceProviderResource: tr,
	}
}

// DiscoverServicesOnce returns a function that guarantees that service discovery only happens once -
// even in concurrent usage.
func DiscoverServicesOnce(project, serviceName string, tracerProviderOptions []sdktrace.TracerProviderOption, opts ...discoveryOption) func() (*Services, error) {
	return sync.OnceValues(func() (*Services, error) {
		return discoverServices(project, serviceName, tracerProviderOptions, opts...)
	})
}

type DiscoverServicesResult struct {
	Services *Services
	Error    error
}

// DiscoverServices starts service discovery asynchronously.
// Once it is done, it returns the discovery result in a channel.
func DiscoverServices(project, serviceName string, tracerProviderOptions []sdktrace.TracerProviderOption, opts ...discoveryOption) chan DiscoverServicesResult {
	resultChan := make(chan DiscoverServicesResult, 1)
	go func() {
		defer close(resultChan)

		s, err := DiscoverServicesOnce(project, serviceName, tracerProviderOptions, opts...)()
		resultChan <- DiscoverServicesResult{Services: s, Error: err}
	}()
	return resultChan
}

// discoverServices builds clients for all Services that we use.
func discoverServices(project, serviceName string, tracerProviderOptions []sdktrace.TracerProviderOption, opts ...discoveryOption) (*Services, error) {
	loggingClient, err := NewLoggingClient(project)
	if err != nil {
		return nil, err
	}
	logger := loggingClient.Logger("ErrorReporting")

	errorReportingCtx := context.Background()
	errorClient, err := errorreporting.NewClient(errorReportingCtx, project, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			logger.Log(logging.Entry{
				Severity: logging.Alert,
				Payload:  fmt.Sprintf("Error reporting failed: %s", err),
			})
			// Ignore err since this probably is due to problems with
			// Permissions or Quotas and thus this needs to be fixed first
		},
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	crd, err := oauth.NewApplicationDefault(
		ctx,
		cloudPlatformScope,
	)
	if err != nil {
		return nil, fmt.Errorf("configure trace credentials: %w", err)
	}

	traceExporterOptions := []texporter.Option{
		texporter.WithEndpointURL(telemetryEndpoint),
		texporter.WithDialOption(grpc.WithPerRPCCredentials(crd)),
	}
	var traceResource *resource.Resource
	for _, opt := range opts {
		if opt.traceExporterOption != nil {
			traceExporterOptions = append(traceExporterOptions, opt.traceExporterOption)
		}
		if opt.traceProviderResource != nil {
			traceResource, err = resource.Merge(traceResource, opt.traceProviderResource)
			if err != nil {
				return nil, fmt.Errorf("merge trace resource: %w", err)
			}

		}
	}

	exporter, err := texporter.New(ctx, traceExporterOptions...)
	if err != nil {
		panic(err)
	}

	s := &Services{
		Logging:        loggingClient,
		ErrorReporting: errorClient,
	}

	discoverPubSub := false
	for _, opt := range opts {
		if opt.gkeAutoDiscoverMetaData {
			metadata, err := gke.GetMetaData(ctx)
			if err != nil {
				logger.Log(logging.Entry{
					Severity: logging.Info,
					Payload:  fmt.Sprintf("Error getting MetaData: %s", err),
				})
				continue
			}

			s.MonitoredResource = gke.MonitoredResourceFromMetaData(metadata)
			traceResource, err = resource.Merge(traceResource, gke.TraceResourceFromMetaData(serviceName, metadata))
			if err != nil {
				return nil, fmt.Errorf("merge trace resource: %w", err)
			}
		} else if opt.pod != "" {
			s.MonitoredResource = gke.MonitoredResource(ctx, s.Logging, project, opt.clusterName, opt.namespace, opt.pod, opt.containerName)
		} else if opt.discoverPubSub {
			discoverPubSub = true
		}
	}

	traceResource, err = resource.Merge(traceResource, resource.NewSchemaless(
		attribute.String("gcp.project_id", project),
		semconv.ServiceName(serviceName),
	))
	if err != nil {
		return nil, fmt.Errorf("merge trace resource: %w", err)
	}

	tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithResource(traceResource))

	if discoverPubSub {
		s.PubSub, err = pubsub.NewClient(ctx, project)
		if err != nil {
			return nil, err
		}
	}

	s.TracerProvider = sdktrace.NewTracerProvider(append(tracerProviderOptions, sdktrace.WithBatcher(exporter))...)

	return s, nil
}

// Close closes all Clients that were created.
// Does **not** handle errors in close since there usually
// is not much that can be done on Close failure anyway.
func (s *Services) Close() {
	if s.PubSub != nil {
		s.PubSub.Close()
	}

	tCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	s.TracerProvider.Shutdown(tCtx) // flushes any pending spans
	s.ErrorReporting.Close()
	s.Logging.Close()
}
