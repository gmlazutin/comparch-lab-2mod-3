package main

import (
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"

	"github.com/gmlazutin/comparch-lab-2mod-3/logging"
)

type ImgProcessorOptions struct {
	Logger *slog.Logger
}

func encodeImage(format string, input *image.RGBA, output io.Writer) error {
	switch format {
	case "png":
		return png.Encode(output, input)
	case "jpeg":
		return jpeg.Encode(output, input, &jpeg.Options{
			Quality: 90,
		})
	default:
		panic("invertImage: unknown format: " + format)
	}
}

func InvertImageProcessor(options ImgProcessorOptions) ImageProcessor {
	return func(ctx context.Context, input io.Reader, output io.Writer) error {
		var logger *slog.Logger
		if options.Logger != nil {
			logger = options.Logger.With(MakeDebugLoggerAttrs(ctx)...)
		} else {
			logger = logging.EmptyLogger()
		}

		img, format, err := image.Decode(input)
		if err != nil {
			return err
		}
		bounds := img.Bounds()

		logger.Debug("processing image...",
			slog.String("format", format),
			slog.Int("BoundsX", bounds.Size().X),
			slog.Int("BoundsX", bounds.Size().Y))

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

		return encodeImage(format, inverted, output)
	}
}
