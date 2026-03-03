package imgpool

import (
	"context"
	"io"
	"log/slog"
)

type ImageWriter interface {
	io.WriteCloser

	Commit() error
}

type ImageCollector func(context.Context, string) (ImageWriter, error)

type ImageCollectorOptions struct {
	Logger *slog.Logger
}
