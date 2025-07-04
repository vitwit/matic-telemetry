package client

import (
	"context"

	"github.com/rs/zerolog"
)

// Base custom context wrapper
type AppContext struct {
	context.Context

	Logger zerolog.Logger
}

// Builder methods
func NewAppContext(base context.Context, logger zerolog.Logger) *AppContext {
	return &AppContext{Context: base, Logger: logger}
}
