package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/internal/app"
	"github.com/niksmo/receipt/pkg/logger"
)

func main() {

	config := config.Load()

	logger := logger.New(config.LogLevel)

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)
	defer cancel()

	opts := app.Opts{
		Addr:     config.Addr,
		Login:    config.Login,
		Password: config.Password,
		SMTPHost: config.SMTPHost,
		SMTPPort: config.SMTPPort,
	}
	app := app.New(logger, opts)
	defer app.Close()

	go app.Run(cancel)

	<-ctx.Done()
}
