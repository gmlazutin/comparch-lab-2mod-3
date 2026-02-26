package collectors

import (
	"bytes"
	"context"
	"sync"

	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool"
)

// implements: ImageWriter
type MemoryImgCollectorWriter struct {
	buf        *bytes.Buffer
	collection *sync.Map //map[string][]byte
	name       string
}

func (d *MemoryImgCollectorWriter) Write(p []byte) (int, error) {
	return d.buf.Write(p)
}

func (d *MemoryImgCollectorWriter) Close() error {
	return nil
}

func (d *MemoryImgCollectorWriter) Commit() error {
	d.collection.Store(d.name, d.buf.Bytes())

	return nil
}

func (d *MemoryImgCollectorWriter) Abort() error {
	return nil
}

func MemoryImgCollector(collection *sync.Map) imgpool.ImageCollector {
	return func(ctx context.Context, name string) (imgpool.ImageWriter, error) {
		m := &MemoryImgCollectorWriter{
			buf:        &bytes.Buffer{},
			collection: collection,
			name:       name,
		}

		return m, nil
	}
}
