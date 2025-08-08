package main

import (
	"net/http"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/internal/receipt_service/adapter"
	"github.com/niksmo/receipt/internal/receipt_service/core/service"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/niksmo/receipt/pkg/middleware"
	"github.com/niksmo/receipt/pkg/sig"
)

func main() {
	PrintAppTitle()
	cfg := config.LoadConfig()
	cfg.Print()

	log := logger.New(cfg.LogLevel)

	sigCtx, stop := sig.NotifyContext()
	defer stop()

	kafkaProducer := adapter.NewKafkaProducer(
		log, cfg.SeedBrokers, cfg.Topic)

	kafkaProducer.InitTopic(sigCtx, cfg.Partitions,
		cfg.ReplicationFactor, OnInitTopicFall(log, stop))

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
