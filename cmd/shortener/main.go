package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	"github.com/VladKvetkin/shortener/internal/app/config"
	"github.com/VladKvetkin/shortener/internal/app/handler"
	"github.com/VladKvetkin/shortener/internal/app/router"
	"github.com/VladKvetkin/shortener/internal/app/server"
	"github.com/VladKvetkin/shortener/internal/app/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %v\n", buildVersion)
	fmt.Printf("Build date: %v\n", buildDate)
	fmt.Printf("Build commit: %v\n", buildCommit)

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	storage, err := storage.GetStorage(config)
	if err != nil {
		panic(err)
	}

	defer storage.Close()

	router := router.NewRouter(handler.NewHandler(storage, config))
	server := server.NewServer(config, router.Router)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		zap.L().Info("Running server", zap.String("Address", config.Address))

		if err = server.Start(); err != nil {
			zap.L().Info("error starting server", zap.Error(err))
			return err
		}

		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := server.Stop(); err != nil {
			zap.L().Info("error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
