package codec

import (
	"context"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool"
)

var (
	ErrFormat = errors.New("imgcodec: unknown format")
)

type Options struct {
	Jpeg *jpeg.Options
}

type imageFormatEncodeFunc func(io.Writer, image.Image, Options) error

var allowedImageFormats = map[string]imageFormatEncodeFunc{
	"png": func(w io.Writer, i image.Image, o Options) error {
		return png.Encode(w, i)
	},
	"jpeg": func(w io.Writer, i image.Image, o Options) error {
		return jpeg.Encode(w, i, o.Jpeg)
	},
}

type codec struct {
	opts Options
}

func New(opts Options) imgpool.Codec {
	return codec{
		opts: opts,
	}
}

func (c codec) Encode(format string, input image.Image, output io.Writer) error {
	if enc, ok := allowedImageFormats[format]; ok {
		return enc(output, input, c.opts)
	}

	return ErrFormat
}

func (c codec) Decode(input io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(input)
	if err != nil {
		if err == image.ErrFormat {
			return nil, "", ErrFormat
		}
		return nil, "", err
	}
	if _, ok := allowedImageFormats[format]; !ok {
		return nil, "", ErrFormat
	}

	return img, format, nil
}

func (c codec) EncodeContext(ctx context.Context, format string, input image.Image, output io.Writer) error {
	return c.Encode(format, input, output)
}

func (c codec) DecodeContext(ctx context.Context, input io.Reader) (image.Image, string, error) {
	return c.Decode(input)
}
