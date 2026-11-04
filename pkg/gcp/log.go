package gcp

import (
	"context"
	"fmt"
	"os"
	"time"

	gcpdetector "go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

type logConfig struct {
	exporter sdklog.Exporter
	resource *resource.Resource
}

type LoggerProviderOption func(*logConfig)

func WithLogExporter(exporter sdklog.Exporter) LoggerProviderOption {
	return func(cfg *logConfig) { cfg.exporter = exporter }
}

func WithLogResource(res *resource.Resource) LoggerProviderOption {
	return func(cfg *logConfig) { cfg.resource = res }
}

func NewLoggerProviderWithProjectExporter(
	ctx context.Context,
	project, serviceName string,
	opts ...LoggerProviderOption,
) (*sdklog.LoggerProvider, error) {
	if project == "" || serviceName == "" {
		return nil, fmt.Errorf("project and service name must not be empty")
	}
	cfg := logConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	detectCtx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	detected, err := resource.New(detectCtx,
		resource.WithDetectors(gcpdetector.NewDetector()),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "detect GCP resource: %v\n", err)
	}

	res, err := resource.Merge(detected, cfg.resource)
	if err != nil {
		return nil, fmt.Errorf("merge log resource: %w", err)
	}
	res = resource.NewWithAttributes(
		res.SchemaURL(),
		append(res.Attributes(),
			attribute.String("gcp.project_id", project),
			semconv.ServiceName(serviceName),
		)...,
	)

	if cfg.exporter == nil {
		auth, err := oauth.NewApplicationDefault(ctx, "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			return nil, fmt.Errorf("configure log credentials: %w", err)
		}
		cfg.exporter, err = otlploggrpc.New(ctx,
			otlploggrpc.WithEndpointURL("https://telemetry.googleapis.com:443"),
			otlploggrpc.WithDialOption(grpc.WithPerRPCCredentials(auth)),
		)
		if err != nil {
			return nil, fmt.Errorf("create log exporter: %w", err)
		}
	}
	return sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(cfg.exporter)),
	), nil
}
