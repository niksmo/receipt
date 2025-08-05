package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/internal/adapter"
	"github.com/niksmo/receipt/internal/core/service"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New(cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)
	defer cancel()

	kafkaProducerOpts := adapter.KafkaProducerOpts{
		SeedBrokers:       cfg.BrokerConfig.SeedBrokers,
		Topic:             cfg.BrokerConfig.Topic,
		Partitions:        cfg.BrokerConfig.Partitions,
		ReplicationFactor: cfg.BrokerConfig.ReplicationFactor,
		MinInsyncReplicas: cfg.BrokerConfig.MinInsyncReplicas,
	}
	kafkaProducer, err := adapter.NewKafkaProducer(ctx, log, kafkaProducerOpts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create kafka producer")
	}

	httpServer := httpserver.New(log, cfg.HTTPServerAddr)
	defer httpServer.Close()

	service := service.NewService(log, kafkaProducer)
	adapter.RegisterMailReceiptHandler(log, httpServer.Mux(), service)
	httpServerStopped := httpServer.Run()

	select {
	case <-ctx.Done():
	case <-httpServerStopped:
	}
}
