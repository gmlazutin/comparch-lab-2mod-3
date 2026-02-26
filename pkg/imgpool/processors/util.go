package processors

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool"
)

func EncodeImage(format string, opts imgpool.ImageProcessorOptions, input *image.RGBA, output io.Writer) error {
	switch format {
	case "png":
		return png.Encode(output, input)
	case "jpeg":
		return jpeg.Encode(output, input, opts.JpegOptions)
	default:
		panic("encodeImage: unknown format: " + format)
	}
}
