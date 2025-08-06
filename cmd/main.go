package main

import (
	"context"
	"os"
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
	cfg.Print(os.Stdout)

	log := logger.New(cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)
	defer cancel()

	kafkaProducer := createKafkaProducer(ctx, log, cfg.BrokerConfig)
	defer kafkaProducer.Close()

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

func createKafkaProducer(
	ctx context.Context, log logger.Logger, brokerCfg config.BrokerConfig,
) *adapter.KafkaProducer {
	kafkaProducerOpts := adapter.KafkaProducerOpts{
		SeedBrokers:       brokerCfg.SeedBrokers,
		Topic:             brokerCfg.Topic,
		Partitions:        brokerCfg.Partitions,
		ReplicationFactor: brokerCfg.ReplicationFactor,
	}
	kafkaProducer, err := adapter.NewKafkaProducer(ctx, log, kafkaProducerOpts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create kafka producer")
	}
	return kafkaProducer
}
