package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/gmlazutin/comparch-lab-2mod-3/pool"
)

type imgPoolCtxKey int

const (
	imgPoolThreadId imgPoolCtxKey = 1
)

func ExtractImagePoolThreadId(ctx context.Context) int {
	return ctx.Value(imgPoolThreadId).(int)
}

type Image struct {
	Name string
	Img  io.Reader
}

type ImageErrorCollector func(context.Context, Image, error)

type ImageWriter interface {
	io.WriteCloser
	Commit() error
	Abort() error
}

type ImageProcessor func(context.Context, io.Reader, io.Writer) error
type ImageCollector func(context.Context, Image) (ImageWriter, error)

type ImagePool struct {
	qpool        *pool.QPool[Image]
	imgprocessor ImageProcessor
	imgcollector ImageCollector
	errcollector ImageErrorCollector
}

func NewImagePool(workers int, processor ImageProcessor, collector ImageCollector) *ImagePool {
	ip := &ImagePool{
		imgprocessor: processor,
		imgcollector: collector,
	}
	ip.qpool = pool.NewQPool(workers, func(ctx context.Context, thread int, input Image) {
		ctx = context.WithValue(ctx, imgPoolThreadId, thread)
		err := ip.process(ctx, input)
		if err != nil && ip.errcollector != nil {
			ip.errcollector(ctx, input, fmt.Errorf("imagepool: %w", err))
		}
	})

	return ip
}

func (ip *ImagePool) WithErrorCollector(collector ImageErrorCollector) *ImagePool {
	ip.errcollector = collector
	return ip
}

func (ip *ImagePool) process(ctx context.Context, input Image) error {
	buf, err := ip.imgcollector(ctx, input)
	if err != nil {
		return fmt.Errorf("error during initializing collector: %w", err)
	}
	defer buf.Close()
	err = ip.imgprocessor(ctx, input.Img, buf)
	if err != nil {
		return fmt.Errorf("image processing fail: %w", errors.Join(err, buf.Abort()))
	}
	if err = buf.Commit(); err != nil {
		return fmt.Errorf("image saving fail: %w", err)
	}

	return nil
}

func (ip *ImagePool) Push(img Image) error {
	return ip.qpool.Push(img)
}

func (ip *ImagePool) PushContext(ctx context.Context, img Image) error {
	return ip.qpool.PushContext(ctx, img)
}

func (ip *ImagePool) Wait() {
	ip.qpool.Wait()
}
