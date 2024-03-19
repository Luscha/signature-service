package logger

import (
	"context"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func GetLogger(ctx context.Context) *Logger {
	r := ctx.Value(loggerKey{})
	if r == nil {
		return Main
	}

	return r.(*Logger)
}

func WithDebugLogger(ctx context.Context) context.Context {
	return WithLogger(ctx, NewLogger("DEBUG", DebugLevel, map[string]interface{}{}))
}
