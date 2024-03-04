//go:build !dev

package opentelemetry

import (
	"context"
	"errors"
	"fmt"
)

type ProdTelemetry struct{}

func (p *ProdTelemetry) Setup(ctx context.Context) (func(context.Context) error, error) {
	fmt.Println("Setting up telemetry for Prod environment")

	shutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup OpenTelemetry SDK failed: %w", err)
	}

	return shutdown, nil
}

func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	return
}
