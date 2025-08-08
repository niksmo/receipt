package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/internal/receipt_service/adapter"
	"github.com/niksmo/receipt/internal/receipt_service/core/service"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/niksmo/receipt/pkg/middleware"
	"github.com/niksmo/receipt/pkg/sig"
)

func main() {
	printAppTitle()
	cfg := config.LoadConfig()
	cfg.Print(os.Stdout)

	log := logger.New(cfg.LogLevel)

	sigCtx, stop := sig.NotifyContext()
	defer stop()

	kafkaProducer := adapter.NewKafkaProducer(
		log, cfg.SeedBrokers, cfg.Topic)

	kafkaProducer.InitTopic(sigCtx, cfg.Partitions,
		cfg.ReplicationFactor, onInitTopicFall(log, stop))

	service := service.NewService(log, kafkaProducer)

	kafkaConsumer := adapter.NewKafkaConsumer(
		log, cfg.SeedBrokers, cfg.Topic, cfg.ConsumerGroup, service)

	mux := http.NewServeMux()
	adapter.RegisterMailReceiptHandler(log, mux, service)

	httpHandler := middleware.LogResposeStatus(log, middleware.AcceptJSON(mux))
	httpServer := httpserver.New(log, cfg.HTTPServerAddr, httpHandler)
	go httpServer.Run(stop)
	go kafkaConsumer.Run(sigCtx)

	<-sigCtx.Done()
	httpServer.Close()
	kafkaProducer.Close()
	kafkaConsumer.Close()
}

func onInitTopicFall(log logger.Logger, stop context.CancelFunc) func(error) {
	return func(err error) {
		log.Error().Err(err).Msg("failed to init broker topic")
		stop()
	}
}

func printAppTitle() {
	fmt.Printf(`
+-----------------------+
|🧾RECEIPT APPLICATION🚀|
+-----------------------+
`)
}
