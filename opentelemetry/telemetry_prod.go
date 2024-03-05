//go:build !dev

package opentelemetry

import (
	"context"
)

type ProdTelemetry struct{}

func NewTelemetryProvider() TelemetryProvider {
	return &ProdTelemetry{}
}

func (p *ProdTelemetry) Setup(_ context.Context) (func(context.Context) error, error) {
	return func(_ context.Context) error { return nil }, nil
}
