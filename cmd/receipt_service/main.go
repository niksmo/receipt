package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/internal/receipt_service/adapter"
	"github.com/niksmo/receipt/internal/receipt_service/core/service"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/niksmo/receipt/pkg/middleware"
)

func main() {
	cfg := config.LoadConfig()
	cfg.Print(os.Stdout)

	log := logger.New(cfg.LogLevel)

	sigCtx, stop := notifyContext()
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
	go httpServer.Run(sigCtx, onHTTPServerFall(log, stop))
	go kafkaConsumer.Run(sigCtx)

	<-sigCtx.Done()
	httpServer.Close()
	kafkaProducer.Close()
	kafkaConsumer.Close()
}

func notifyContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
}

func onInitTopicFall(log logger.Logger, stop context.CancelFunc) func(error) {
	return func(err error) {
		log.Error().Err(err).Msg("failed to init broker topic")
		stop()
	}
}

func onHTTPServerFall(log logger.Logger, stop context.CancelFunc) func(error) {
	return func(err error) {
		log.Error().Err(err).Msg("http server crashed")
		stop()
	}
}
