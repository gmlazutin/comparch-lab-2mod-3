package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/gmlazutin/comparch-lab-2mod-3/logging"
	"github.com/gmlazutin/comparch-lab-2mod-3/util"
)

func main() {
	log_lvl := flag.Int("log", 0, "log level")
	imgs_dir := flag.String("input", "./input", "source directory")
	output_dir := flag.String("output", "./output", "destination directory")
	workers_cnt := flag.Int("workers", 4, "workers count, must be in (0, 500]")
	algo := flag.String("algo", "invert", "needed action, only \"invert\" action is currently supported")
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
	pool := NewImagePool(
		stopctx,
		*workers_cnt,
		InvertImageProcessor(ImgProcessorOptions{
			Logger: logger,
		}),
		MemoryImgCollector(collection),
	).WithErrorCollector(func(ctx context.Context, i Image, err error) {
		logger.Error(
			"image pool error has occurred",
			append(MakeDebugLoggerAttrs(ctx), logging.Error(err))...,
		)
	})

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			logger.Error("unable to open file", logging.Error(err))
			continue
		}

		err = pool.PushContext(stopctx, Image{
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
