package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/receipt/internal/receiver/adapters"
	"github.com/niksmo/receipt/pkg/logger"
)

var (
	seedBrokers = []string{
		"127.0.0.1:19094",
		"127.0.0.1:29094",
		"127.0.0.1:29094",
	}
	topic    = "mail-receipt"
	logLevel = "debug"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer cancel()

	logger := logger.New(logLevel)

	producer, err := adapters.NewMessageProducer(
		ctx, logger, seedBrokers, topic,
	)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}
	defer producer.Close()
}
