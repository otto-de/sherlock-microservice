package gcp

import (
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newTraceProviderWithProjectExporter(project string, opts ...sdktrace.TracerProviderOption) *sdktrace.TracerProvider {
	exporter, err := texporter.New(texporter.WithProjectID(project))
	if err != nil {
		panic(err)
	}
	return sdktrace.NewTracerProvider(append(opts, sdktrace.WithBatcher(exporter))...)
}

func MustInitTracerProvider(project string, opts ...sdktrace.TracerProviderOption) *sdktrace.TracerProvider {
	tp := newTraceProviderWithProjectExporter(project, opts...)
	otel.SetTracerProvider(tp)
	return tp
}
