package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"taskhub/config"
	"taskhub/internal/gateway"
	"taskhub/pkg/logger"
)

func main() {
	ctx := context.Background()
	config := config.NewConfig()
	logger := logger.NewLogger()

	gw := gateway.NewGateway(config, logger)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		if err := gw.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown: %v", err)
		}
	}()

	if err := gw.Start(); err != nil {
		logger.Error("Failed to start: %v", err)
	}
}
