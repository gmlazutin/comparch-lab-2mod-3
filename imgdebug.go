package main

import (
	"context"
	"log/slog"
)

func MakeDebugLoggerAttrs(ctx context.Context) []any {
	return []any{
		slog.Int("thread", ExtractImagePoolThreadId(ctx)),
		slog.String("element", ExtractImagePoolName(ctx)),
	}
}
