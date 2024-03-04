package opentelemetry

import "context"

type TelemetryProvider interface {
	Setup(ctx context.Context) (shutdown func(context.Context) error, err error)
}
