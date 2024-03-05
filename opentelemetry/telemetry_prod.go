//go:build !dev

package opentelemetry

import (
	"context"
)

type ProdTelemetry struct{}

func (p *ProdTelemetry) Setup(ctx context.Context) (func(context.Context) error, error) {
	return func(ctx.Context) error { return nil }, nil
}
