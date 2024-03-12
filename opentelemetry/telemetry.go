package opentelemetry

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	once sync.Once
)

func Setup(ctx context.Context) (func(context.Context) error, error) {
	var shutdown func(context.Context) error
	var err error

	once.Do(func() {
		fmt.Println("Setting up telemetry for DEV environment")

		prop := newPropagator()
		otel.SetTextMapPropagator(prop)

		tracerProvider, innerErr := newTraceProvider(ctx)
		if innerErr != nil {
			err = fmt.Errorf("setup OpenTelemetry SDK failed: %w", innerErr)
			return
		}
		otel.SetTracerProvider(tracerProvider)

		shutdown = func(ctx context.Context) error {
			if innerErr := tracerProvider.Shutdown(ctx); innerErr != nil {
				return fmt.Errorf("tracerProvider shutdown failed: %w", innerErr)
			}
			return nil
		}
		fmt.Println("Telemetry setup completed....")
	})

	if err != nil {
		return nil, err
	}
	return shutdown, nil
}

func Shutdown(ctx context.Context, shutdown func(context.Context) error) {
	if shutdownErr := shutdown(ctx); shutdownErr != nil {
		fmt.Printf("Error during OpenTelemetry shutdown: %v\n", shutdownErr)
	}
}

// GetTracer You must use it after calling the Setup() function.
func GetTracer(tracerName string) trace.Tracer {
	if tracerName == "" {
		tracerName = "default_tracer"
	}

	return otel.Tracer(tracerName)
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
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

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return traceProvider, nil
}
