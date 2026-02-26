package main

import (
	"context"
	"flag"
	"image/jpeg"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/gmlazutin/comparch-lab-2mod-3/internal/logging"
	"github.com/gmlazutin/comparch-lab-2mod-3/internal/util"
	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool"
	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool/collectors"
	"github.com/gmlazutin/comparch-lab-2mod-3/pkg/imgpool/processors"
)

func main() {
	log_lvl := flag.Int("log", 0, "log level")
	imgs_dir := flag.String("input", "./input", "source directory")
	output_dir := flag.String("output", "./output", "destination directory")
	workers_cnt := flag.Int("workers", 4, "workers count, must be in (0, 500]")
	algo := flag.String("algo", "invert", "needed action, only \"invert\" action is currently supported")
	timing := flag.Bool("timing", false, "-timing")
	flag.Parse()

	logger := logging.InitLogger(slog.Level(*log_lvl))

	if _, err := os.Stat(*imgs_dir); os.IsNotExist(err) {
		logger.Error("input directory does not exist", slog.String("dest", *imgs_dir))
		return
	}

	if err := os.MkdirAll(*output_dir, 0755); err != nil {
		logger.Error("failed to create output directory", slog.String("dest", *output_dir), logging.Error(err))
		return
	}

	if 0 >= *workers_cnt || *workers_cnt > 500 {
		logger.Error("invalid workers count value", slog.Int("value", *workers_cnt))
		return
	}

	if *algo != "invert" {
		logger.Error("unsupported action", slog.String("algo", *algo))
		return
	}

	files, err := util.ListFilesWithExts(*imgs_dir, []string{".jpg", ".jpeg", ".png"})
	if err != nil {
		logger.Error("unable to get input directory listing", logging.Error(err))
		return
	}

	stopctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	collection := &sync.Map{}
	pool := imgpool.NewImagePool(
		stopctx,
		*workers_cnt,
		processors.InvertImageProcessor(imgpool.ImageProcessorOptions{
			Logger: logger,
			JpegOptions: &jpeg.Options{
				Quality: 70,
			},
		}),
		collectors.MemoryImgCollector(collection),
	).WithErrorCollector(func(ctx context.Context, i imgpool.Image, err error) {
		logger.Error(
			"image pool error has occurred",
			append(imgpool.MakeDebugLoggerAttrs(ctx), logging.Error(err))...,
		)
	})

	start := time.Now()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			logger.Error("unable to open file", logging.Error(err))
			continue
		}

		err = pool.PushContext(stopctx, imgpool.Image{
			Name: f.Name(),
			Img:  f,
		})

		if err != nil {
			break
		}
	}

	err = pool.WaitDone()
	if err != nil {
		logger.Info("stopping processing...")
		return
	}
	logger.Info("writing output...")

	if *timing {
		logger.Info("done", slog.Duration("time passed", time.Since(start)))
	}

	collection.Range(func(key, value any) bool {
		if stopctx.Err() != nil {
			logger.Info("stopping saving...")
			return false
		}
		if err := os.WriteFile(filepath.Join(*output_dir, filepath.Base(key.(string))), value.([]byte), 0600); err != nil {
			logger.Error("unable to save output file", logging.Error(err))
		}
		return true
	})
}
