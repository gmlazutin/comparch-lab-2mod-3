package main

import (
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
)

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

func InvertImage(ctx context.Context, input io.Reader, output io.Writer) error {
	img, format, err := image.Decode(input)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
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

	return encodeImage(format, inverted, output)
}
