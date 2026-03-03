package imgpool

import (
	"context"
	"image"
	"io"
)

type Codec interface {
	EncodeContext(ctx context.Context, format string, input image.Image, output io.Writer) error
	Encode(format string, input image.Image, output io.Writer) error
	DecodeContext(ctx context.Context, input io.Reader) (image.Image, string, error)
	Decode(input io.Reader) (image.Image, string, error)
}
