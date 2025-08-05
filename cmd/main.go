package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New(cfg.LogLevel)

	httpServer := httpserver.New(log, cfg.HTTPServerAddr)
	defer httpServer.Close()

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)
	defer cancel()

	httpServerStopped := httpServer.Run()

	select {
	case <-ctx.Done():
	case <-httpServerStopped:
	}
}
