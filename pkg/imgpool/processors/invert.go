package processors

import (
	"context"
	"image"
	"image/color"
	"io"
	"log/slog"

	"github.com/gmlazutin/comparch-lab-2mod-3/internal/logging"
	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool"
)

func InvertImageProcessor(options imgpool.ImageProcessorOptions) imgpool.ImageProcessor {
	return func(ctx context.Context, input io.Reader, output io.Writer) error {
		var logger *slog.Logger
		if options.Logger != nil {
			logger = options.Logger.With(imgpool.MakeDebugLoggerAttrs(ctx)...)
		} else {
			logger = logging.EmptyLogger()
		}

		img, format, err := image.Decode(input)
		if err != nil {
			return err
		}
		bounds := img.Bounds()

		sz := bounds.Size()
		logger.Debug("processing image...",
			slog.String("format", format),
			slog.Int("BoundsX", sz.X),
			slog.Int("BoundsY", sz.Y))

		inverted := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				inverted.Set(x, y, color.RGBA{
					uint8(255 - r/256),
					uint8(255 - g/256),
					uint8(255 - b/256),
					uint8(a / 256),
				})
				if ctx.Err() != nil {
					return ctx.Err()
				}
			}
		}

		logger.Debug("inverting done, trying to encode...")

		return EncodeImage(format, options, inverted, output)
	}
}
