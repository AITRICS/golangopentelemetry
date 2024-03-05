//go:build !dev

package opentelemetry

import (
	"context"
	"errors"
	"fmt"
)

type ProdTelemetry struct{}

func (p *ProdTelemetry) Setup(ctx context.Context) (func(context.Context) error, error) {
	return func(context.Context) error { return nil }, nil
}
