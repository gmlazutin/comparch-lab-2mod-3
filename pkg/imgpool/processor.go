package imgpool

import (
	"context"
	"io"
	"log/slog"
)

type ImageProcessor func(context.Context, io.Reader, io.Writer) error

type ImageProcessorOptions struct {
	Logger *slog.Logger
	Codec  Codec
}
