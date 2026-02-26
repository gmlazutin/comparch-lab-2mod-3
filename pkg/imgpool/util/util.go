package util

import (
	"fmt"
	"os"

	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool"
)

type ErrFileSizeLimitExceeded struct {
	Path   string
	Actual int64
	Needed int64
}

func (e ErrFileSizeLimitExceeded) Error() string {
	return fmt.Sprintf("size limit exceeded for %q (%d bytes vs. %d bytes max)", e.Path, e.Actual, e.Needed)
}

const (
	DEFAULT_MAX_IMG_SIZE = 1024 * 1024 * 50 //50mb
)

func OpenImage(name string, limit int64, out *imgpool.Image) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("openImage: unable to open file: %w", err)
	}
	fstat, err := f.Stat()
	if err != nil {
		f.Close()
		return fmt.Errorf("openImage: unable to fstat file: %w", err)
	}
	if fstat.Size() > limit {
		f.Close()
		return fmt.Errorf("openImage: %w", ErrFileSizeLimitExceeded{
			Path:   name,
			Actual: fstat.Size(),
			Needed: limit,
		})
	}

	out.Img = f
	out.Name = name

	return nil
}
