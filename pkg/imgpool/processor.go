package imgpool

import (
	"log/slog"
)

type ImageProcessorOptions struct {
	Logger *slog.Logger
	Codec  Codec
}
