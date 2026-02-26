package imgpool

import (
	"image/jpeg"
	"log/slog"
)

type ImageProcessorOptions struct {
	Logger      *slog.Logger
	JpegOptions *jpeg.Options
}
