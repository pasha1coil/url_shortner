package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"urlshortner/internal/app"
	"urlshortner/internal/initialize"
)

// run cmd/main.go --in_memory

func main() {
	use := flag.Bool("in_memory", false, "in_memory")
	flag.Parse()
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	config, err := initialize.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = app.Run(ctx, config, logger, *use); err != nil {
		logger.Fatal("App exited with error", zap.Error(err))
	}
}
