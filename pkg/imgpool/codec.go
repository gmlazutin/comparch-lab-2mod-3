package imgpool

import (
	"image"
	"io"
)

type Codec interface {
	Encode(format string, input image.Image, output io.Writer) error
	Decode(input io.Reader) (image.Image, string, error)
}
