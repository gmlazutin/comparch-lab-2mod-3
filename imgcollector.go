package main

import (
	"bytes"
	"context"
	"sync"
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

func MemoryImgCollector(collection *sync.Map) ImageCollector {
	return func(ctx context.Context, img Image) (ImageWriter, error) {
		m := &MemoryImgCollectorWriter{
			buf:        &bytes.Buffer{},
			collection: collection,
			name:       img.Name,
		}

		return m, nil
	}
}
