//go:build !dev

package opentelemetry

import (
	"context"
)

type ProdTelemetry struct{}

func (p *ProdTelemetry) Setup(_ context.Context) (func(context.Context) error, error) {
	return func(_ context.Context) error { return temp }, nil
}
