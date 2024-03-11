package opentelemetry

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Setup(ctx context.Context) (func(context.Context) error, error) {
	fmt.Println("Setting up telemetry for DEV environment")

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup OpenTelemetry SDK failed: %w", err)
	}
	otel.SetTracerProvider(tracerProvider)

	shutdown := func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("tracerProvider shutdown failed: %w", err)
		}
		return nil
	}

	return shutdown, nil
}

func Shutdown(shutdown func(context.Context) error, ctx context.Context) {
	if shutdownErr := shutdown(ctx); shutdownErr != nil {
		fmt.Printf("Error during OpenTelemetry shutdown: %v\n", shutdownErr)
	}
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter failed: %w", err)
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "default_service_name"
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource failed: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	return traceProvider, nil
}
