package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"backend/internal/config"
	"backend/internal/http/handler"
	"backend/internal/http/server"
	"backend/internal/logger"
	"backend/internal/repository"
	"backend/internal/usecase"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse()

	cfg := config.ParseConfig(*configPath)

	// fmt.Println(*cfg)

	logger.MustInit(cfg.Logger.Level)

	pg, err := repository.NewConnection(context.Background(), cfg.DB)
	if err != nil {
		logger.Error("error connecting to database", slog.String("error", err.Error()))

		os.Exit(1)
	}

	repo := repository.New(pg)

	pingReportService := usecase.New(repo)

	hnd := handler.New(pingReportService)

	app := server.New(hnd.Route(), cfg.Server)

	go func() {
		if err = app.Run(); err != nil {
			logger.Error("error running server", slog.String("error", err.Error()))
		}
	}()

	logger.Info("starting server", slog.String("address", cfg.Server.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = app.Shutdown(ctx); err != nil {
		logger.Error("shutdown server error", slog.String("error", err.Error()))
	}

}
