package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
	"urlshortner/internal/controller"
	"urlshortner/internal/initialize"
	"urlshortner/internal/repository"
	httpserver "urlshortner/internal/server/http"
	"urlshortner/internal/service"
	"urlshortner/pkg/closer"
)

func Run(ctx context.Context, config *initialize.Config, logger *zap.Logger, use bool) error {
	println(use)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from a panic", zap.Any("error", r))
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	shutdownGroup := closer.NewCloserGroup()

	var (
		err          error
		pgDB         *initialize.DB
		shortnerRepo service.SplitRepository
	)

	if !use {
		pgDB, err = initialize.InitDB(ctx, config)
		if err != nil {
			logger.Error("Error init pg db", zap.Error(err))
			return err
		}
		shutdownGroup.Add(closer.CloserFunc(pgDB.DB.Close))
	}

	redisClient, err := initialize.Redis(ctx, config)
	if err != nil {
		logger.Error("Error init redis client", zap.Error(err))
		return err
	}

	if use {
		shortnerRepo = repository.NewInMemoryRepo()
	} else {
		shortnerRepo = repository.NewShortnerRepo(pgDB.DB)
	}

	shortnerService := service.NewShortnerService(service.Deps{
		Repo:        shortnerRepo,
		Logger:      logger,
		Config:      config,
		RedisClient: redisClient,
	})

	shortenController := controller.NewShortenController(shortnerService, logger)

	server := httpserver.NewServer(httpserver.ServerConfig{
		Logger:      logger,
		Controllers: []httpserver.Controller{shortenController},
	})

	go func() {
		if err := server.Start(config.HTTPHost + ":" + config.HTTPPort); err != nil {
			logger.Error("Server startup error", zap.Error(err))
			cancel()
		}
	}()

	server.ListRoutes()

	shutdownGroup.Add(closer.CloserFunc(server.Shutdown))
	shutdownGroup.Add(closer.CloserFunc(redisClient.Close))

	<-ctx.Done()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutCancel()
	if err := shutdownGroup.Call(timeoutCtx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("Shutdown timed out", zap.Error(err))
		} else {
			logger.Error("Failed to shutdown services gracefully", zap.Error(err))
		}
		return err
	}

	logger.Info("Application has stopped")
	return nil
}
